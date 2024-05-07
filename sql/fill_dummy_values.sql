-- Inserting data into the Brand table
INSERT INTO Brand (Name, Description) VALUES ('Apple', 'American multinational technology company');
INSERT INTO Brand (Name, Description) VALUES ('Samsung', 'South Korean multinational electronics company');

-- Inserting data into the Category table
INSERT INTO Category (Name, Description) VALUES ('Smartphones', 'Mobile devices with advanced computing capabilities');
INSERT INTO Category (Name, Description) VALUES ('Laptops', 'Portable personal computers');

-- Inserting data into the Product table
INSERT INTO Product (ID, Name, BrandName, CategoryName, Description, QtyInStock, Price) VALUES (UUID(), 'iPhone 13', 'Apple', 'Smartphones', 'Latest iPhone model', 100, 999.99);
INSERT INTO Product (ID, Name, BrandName, CategoryName, Description, QtyInStock, Price) VALUES (UUID(), 'Galaxy S21', 'Samsung', 'Smartphones', 'Latest Samsung Galaxy model', 100, 899.99);
INSERT INTO Product (ID, Name, BrandName, CategoryName, Description, QtyInStock, Price) VALUES (UUID(), 'MacBook Pro', 'Apple', 'Laptops', 'Latest MacBook model', 100, 1299.99);
INSERT INTO Product (ID, Name, BrandName, CategoryName, Description, QtyInStock, Price) VALUES (UUID(), 'Galaxy Book Pro', 'Samsung', 'Laptops', 'Latest Samsung laptop model', 100, 1199.99);
