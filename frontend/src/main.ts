type ElementConfig = {
    id: string
    type?: new () => HTMLElement
    required?: boolean
}


function ensureElements(config: ElementConfig[]) {
    const results: Record<string, HTMLElement | null> = {}
    for (const { id, type = HTMLElement, required = true } of config) {

        const element = document.getElementById(id)

        if (!element) {
            if (required) throw new Error(`Missing required element #${id}`)
            results[id] = null
            continue
        }

        if (!(element as any instanceof type)) {
            throw new Error(`Element #${id} is not ${type.name}`);
        }

        results[id] = element
    }
    return results as Record<string, HTMLElement>
}

const ws = new WebSocket("ws://192.168.1.115:8080/api/chat")

function connect() {


    try {
        ws.onopen = () => {
            console.log("Connected to WebSocket server")
        }

        ws.onmessage = (event) => {
            const messages = ensureElements([{ id: "messages", type: HTMLDivElement }])
            const messageDisplay = messages["messages"] as HTMLDivElement
            messageDisplay.innerHTML += `<p>${event.data}</p>`
        }

        ws.onclose = () => {
            console.log("Websocket connection closed, Retrying...")
            setTimeout(connect, 1000)
        }

        ws.onerror = (error) => {
            console.error("WebSocket error:", error);
        };



    } catch (err) {
        if (err instanceof Error) {
            console.error("error : ", err)
        } else {
            console.error("Unknown error")
        }

    }

}

function sendMessage() {
    const input = ensureElements([{ id: "messageInput", type: HTMLInputElement }])
    const messageInput = input["messageInput"] as HTMLInputElement
    let message = messageInput.value;
    ws.send(message);
    messageInput.value = "";
}

function button(buttonId: string, buttonHandler: any) {
    const button = ensureElements([{ id: buttonId, type: HTMLButtonElement }])
    const messageButton = button[buttonId] as HTMLButtonElement
    messageButton.addEventListener("click", buttonHandler)
}




connect()
button("messageButton", sendMessage)
