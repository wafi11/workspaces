create table users (
    id UUID
    username VARCHAR(100),
    email VARCHAR(100),
    password TEXT,
    terminal_url TEXT,
    created_at TIMESTAMP
    updated_at TIMESTAMP
)


create table sessions (
    id UUID primary key,
    user_id UUID,
    refresh_token TEXT,
    ip_address VARCHAR(100),
    user_agent VARCHAR(200),
    

)