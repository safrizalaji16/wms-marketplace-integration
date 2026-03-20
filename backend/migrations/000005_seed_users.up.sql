INSERT INTO users (email, password, role)
VALUES
    (
        'superadmin@mail.com',
        crypt('123456', gen_salt('bf')),
        'superAdmin'
    ),
    (
        'warehouseadmin@mail.com',
        crypt('123456', gen_salt('bf')),
        'warehouseAdmin'
    ),
    (
        'picker@mail.com',
        crypt('123456', gen_salt('bf')),
        'picker'
    ),
    (
        'packer@mail.com',
        crypt('123456', gen_salt('bf')),
        'packer'
    )
ON CONFLICT (email) DO NOTHING;
