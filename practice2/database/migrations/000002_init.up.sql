create table if not exists users (
    id serial primary key,
    name varchar(255) not null,
    email varchar(255) not null,
    age int not null,
    created_at timestamp default now()
);

insert into users (name, email, age) values 
('John Pork', 'john@example.com', 30),
('Aibek Kaibulla', 'aibek@example.com', 19),
('Maison Margiela', 'maison@example.com', 34),
('AXAXXAXXAX BCBCBCBCBCB', 'dsadadsdasd@example.com', 21),
('LOXLXOLXOasd', 'LOXLXOLXOasd@example.com', 21)

