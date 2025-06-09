import { cfg } from './config.ts'

type ElementConfig = {
  id: string
  type?: new () => HTMLElement
  required?: boolean
}

type LoginResponse = {
  message: string
  token: string
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

    if (!((element as any) instanceof type)) {
      throw new Error(`Element #${id} is not ${type.name}`)
    }

    results[id] = element
  }
  return results as Record<string, HTMLElement>
}

var ws = new WebSocket(cfg.WS_BACKEND)
const loginUrl = cfg.HTTP_BACKEND

// let loginResponse:LoginResponse
// Store the token globally
let authToken: string | null = null

function connect() {
  if (!authToken) {
    console.error("Cannot connect: No authentication token available")
    return
  }

  // Close existing connection if any
  if (ws && ws.readyState !== WebSocket.CLOSED) {
    ws.close()
  }

  const url = new URL(cfg.WS_BACKEND)
  url.searchParams.append("token", authToken)
  ws = new WebSocket(url.toString())

  try {
    ws.onopen = () => {
      console.log("Authenticated connection to WebSocket server")
    }

    ws.onmessage = (event) => {
      const messages = ensureElements([
        { id: "messages", type: HTMLDivElement },
      ])
      const messageDisplay = messages["messages"] as HTMLDivElement
      messageDisplay.innerHTML += `<p>${event.data}</p>`
    }

    ws.onclose = () => {
      console.log("Other Websocket connection closed")
      // setTimeout(connect, 1000)
    }

    ws.onerror = (error) => {
      console.error("WebSocket error:", error)
    }
  } catch (err) {
    if (err instanceof Error) {
      console.error("error : ", err)
    } else {
      console.error("Unknown error")
    }
  }
}

function sendMessage() {
  if (!authToken) {
    console.error("Cannot send message: Not authenticated")
    return
  }
  const input = ensureElements([
    { id: "messageInput", type: HTMLInputElement },
  ])
  const messageInput = input["messageInput"] as HTMLInputElement
  let message = messageInput.value
  ws.send(message)
  messageInput.value = ""
}

async function login() {
  var input = ensureElements([{ id: "emailInput", type: HTMLInputElement }])

  const emailInput = input["emailInput"] as HTMLInputElement

  input = ensureElements([{ id: "passwordInput", type: HTMLInputElement }])

  const passwordInput = input["passwordInput"] as HTMLInputElement

  let email = emailInput.value
  let password = passwordInput.value

  try {
    const response = await fetch(loginUrl, {
      method: "POST",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        email: email,
        password: password,
      }),
    })
    // Check for HTTP errors
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    const data = await response.json()
    console.log("Login successful:", data)

    // Clear inputs after successful request
    emailInput.value = ""
    passwordInput.value = ""
    authToken = (data as LoginResponse).token
    connect()
  } catch (error) {
    console.error("Login failed:", error)
    // Handle errors (e.g., show error message to user)
  }
}

function button(buttonId: string, buttonHandler: any) {
  const button = ensureElements([{ id: buttonId, type: HTMLButtonElement }])
  const messageButton = button[buttonId] as HTMLButtonElement
  messageButton.addEventListener("click", buttonHandler)
}

// connect()
button("messageButton", sendMessage)
button("loginButton", login)
