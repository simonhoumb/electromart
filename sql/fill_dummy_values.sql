-- Inserting data into the Brand table
SET @appleUUID = UUID();
SET @samsungUUID = UUID();
INSERT INTO Brand (ID, Name, Description) VALUES (@appleUUID, 'Apple', 'American multinational technology company');
INSERT INTO Brand (ID, Name, Description) VALUES (@samsungUUID, 'Samsung', 'South Korean multinational electronics company');

-- Inserting data into the Category table
SET @smartphonesUUID = UUID();
SET @laptopsUUID = UUID();
INSERT INTO Category (ID, Name, Description) VALUES (@smartphonesUUID, 'Smartphones', 'Mobile devices with advanced computing capabilities');
INSERT INTO Category (ID, Name, Description) VALUES (@laptopsUUID, 'Laptops', 'Portable personal computers');

-- Inserting data into the Product table
INSERT INTO Product (ID, Name, BrandID, CategoryID, Description, QtyInStock, Price) VALUES (UUID(), 'iPhone 13', @appleUUID, @smartphonesUUID, 'Latest iPhone model', 100, 999.99);
INSERT INTO Product (ID, Name, BrandID, CategoryID, Description, QtyInStock, Price) VALUES (UUID(), 'Galaxy S21', @samsungUUID, @smartphonesUUID, 'Latest Samsung Galaxy model', 100, 899.99);
INSERT INTO Product (ID, Name, BrandID, CategoryID, Description, QtyInStock, Price) VALUES (UUID(), 'MacBook Pro', @appleUUID, @laptopsUUID, 'Latest MacBook model', 100, 1299.99);
INSERT INTO Product (ID, Name, BrandID, CategoryID, Description, QtyInStock, Price) VALUES (UUID(), 'Galaxy Book Pro', @samsungUUID, @laptopsUUID, 'Latest Samsung laptop model', 100, 1199.99);
