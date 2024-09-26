INSERT INTO users (name, ip, id, email, created_at, id_int) VALUES 
    ('vasya', '228.223.221.224', ('00000000-0000-0000-0000-' || LPAD('2', 12, '0'))::uuid, 'vasyan@example.com', NOW(), 1),
    ('maxim', '228.223.221.223',('00000000-0000-0000-0000-' || LPAD('1', 12, '0'))::uuid, 'maxim_sonselinexample.com', NOW(), 2),
    ('fanat', '228.223.221.222', ('00000000-0000-0000-0000-' || LPAD('3', 12, '0'))::uuid, '5252@example.com', NOW(), 3),
    ('HR', '228.223.221.221', ('00000000-0000-0000-0000-' || LPAD('4', 12, '0'))::uuid, 'fanathrov@ya.ru', NOW(), 4);