DELETE FROM users
WHERE email IN (
    'superadmin@mail.com',
    'warehouseadmin@mail.com',
    'picker@mail.com',
    'packer@mail.com'
);
