TRUNCATE TABLE users RESTART IDENTITY CASCADE;

INSERT INTO users (name,email,password,role)
VALUES 
    ('admin','admin@admin.com','$2a$10$jEpH2uWzBYHaXnEdTgCp2OQsrICXeg.UljLsgV5n3c2Q9lTFa1CHO','admin');
