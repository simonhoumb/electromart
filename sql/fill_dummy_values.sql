-- Inserting data into the Brand table
INSERT IGNORE INTO Brand (Name, Description) VALUES
                                          ('Apple', 'American multinational technology company'),
                                          ('Samsung', 'South Korean multinational electronics company'),
                                          ('Sony', 'Japanese multinational conglomerate corporation'),
                                          ('LG', 'South Korean multinational electronics company'),
                                          ('Netgear', 'American multinational computer networking company'),
                                          ('Microsoft', 'American multinational technology corporation'),
                                          ('Google', 'American multinational technology company'),
                                          ('HP', 'American multinational information technology company'),
                                          ('Asus', 'Taiwanese multinational computer hardware and consumer electronics company'),
                                          ('Lenovo', 'Chinese multinational technology company'),
                                          ('Acer', 'Taiwanese multinational electronics company'),
                                          ('Razer', 'American multinational computer hardware, software, and peripherals company (gaming focus)');

-- Inserting data into the Category table
INSERT IGNORE INTO Category (Name, Description) VALUES
                                             ('Smartphones', 'Mobile devices with advanced computing capabilities'),
                                             ('Laptops', 'Portable personal computers'),
                                             ('Tablets', 'Mobile computing devices with a touchscreen interface'),
                                             ('Phones', 'Telecommunication devices for voice communication'),
                                             ('Desktops', 'Personal computers designed for stationary use'),
                                             ('Components', 'Individual parts that make up a computer system'),
                                             (  'Audio', 'Electronic devices for reproducing sound'),
                                             (  'Video', 'Electronic devices for displaying visual information'),
                                             (  'Gaming', 'Products designed for video game playing'),
                                             (  'Networking', 'Equipment for connecting computers and devices together'),
                                             (  'Software', 'Applications and programs used with a computer'),
                                             (  'Accessories', 'Additional items that complement electronic devices'),
                                             ('Smart Home', 'Networked devices that can be controlled remotely for home automation'),
                                             ('Appliances', 'Electrical devices used for performing household tasks'),
                                             ('Wearables', 'Electronic devices worn on the body'),
                                             ('Office Electronics', 'Equipment used for office tasks, such as printing and scanning');


-- Inserting data into the Product table
INSERT INTO Product (ID, Name, BrandName, CategoryName, Description, QtyInStock, Price)
VALUES
    (UUID(), 'iPhone 13', 'Apple', 'Smartphones', 'Latest iPhone model', 100, 9999.99),
    (UUID(), 'Galaxy S21', 'Samsung', 'Smartphones', 'Latest Samsung Galaxy model', 100, 8999.99),
    (UUID(), 'MacBook Pro', 'Apple', 'Laptops', 'Latest MacBook model', 100, 1299.99),
    (UUID(), 'Galaxy Book Pro', 'Samsung', 'Laptops', 'Latest Samsung laptop model', 100, 11999.99),
    (UUID(), 'AirPods Pro', 'Apple', 'Audio', 'Wireless earbuds with noise cancellation', 50, 2499.99),
    (UUID(), 'WH-1000XM4', 'Sony', 'Audio', 'Over-ear headphones with superior noise cancellation', 75, 3499.99),
    (UUID(), 'Smart TV', 'LG', 'Video', 'High-definition television with smart features', 20, 5999.99),
    (UUID(), 'Playstation 5', 'Sony', 'Gaming', 'Next-generation video game console', 25, 4999.99),
    (UUID(), 'Wireless Router', 'Netgear', 'Networking', 'High-speed wireless router for home networks', 30, 999.99),
    (UUID(), 'Microsoft Office Suite', 'Microsoft', 'Software', 'Productivity software for home and business use', 100, 1499.99),
    (UUID(), 'Google Nest', 'Google', 'Smart Home', 'Voice-activated speaker for smart home control', 40, 999.99),
    (UUID(), 'Lenovo Tab M10', 'Lenovo', 'Tablets', 'Android tablet for entertainment and productivity', 60, 2999.99)
ON DUPLICATE KEY UPDATE ID  = UUID();