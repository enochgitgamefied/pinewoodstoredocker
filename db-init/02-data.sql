-- Sample products data
INSERT INTO products (name, description, price, is_listed) VALUES
('iPhone 15 Pro', '6.1-inch Super Retina XDR display, A17 Pro chip, 48MP camera', 999.00, TRUE),
('Samsung Galaxy S23 Ultra', '6.8-inch Dynamic AMOLED, 200MP camera, S Pen included', 1199.99, TRUE),
-- [Rest of your product data...]
('Amazon Echo Dot (3rd Gen)', 'Older model with less features', 39.99, FALSE);

-- Sample admin user (password should be hashed in production)
INSERT INTO users (username, password, role) VALUES
('admin', '$2a$10$somehashedpassword', 'ADMIN');