    CREATE TABLE Brand (
                           Name        varchar(128) NOT NULL,
                           Description varchar(512),
                           PRIMARY KEY (Name));
    CREATE TABLE Product (
                             ID           varchar(64) NOT NULL,
                             Name         varchar(128) NOT NULL,
                             BrandName    varchar(128),
                             CategoryName varchar(128) DEFAULT 'Uncategorised' NOT NULL,
                             Description  varchar(512),
                             QtyInStock   int(10) NOT NULL,
                             Price        double NOT NULL,
                             Active       bit(1) DEFAULT 0 NOT NULL,
                             PRIMARY KEY (ID),
                             INDEX (Name),
                             INDEX (BrandName),
                             INDEX (CategoryName));
    CREATE TABLE ProductOrder (
                                  ID                varchar(64) NOT NULL,
                                  UserAccountID     varchar(64) NOT NULL,
                                  OrderDate         date NOT NULL,
                                  ShippedDate       date,
                                  EstimatedDelivery date,
                                  DeliveryFee       double NOT NULL,
                                  DeliveryAddress   varchar(128) NOT NULL,
                                  PostalCode        int(4) NOT NULL,
                                  Status            varchar(64) NOT NULL,
                                  PaymentMethod     varchar(64) NOT NULL,
                                  Comments          varchar(512),
                                  PRIMARY KEY (ID));
    CREATE TABLE UserAccount (
                                 ID        varchar(64) NOT NULL,
                                 Username  varchar(64) NOT NULL,
                                 Email     varchar(128) NOT NULL,
                                 Password  varchar(64) NOT NULL,
                                 FirstName varchar(128) NOT NULL,
                                 LastName  varchar(128) NOT NULL,
                                 Phone     int(12) NOT NULL,
                                 PRIMARY KEY (ID));
    CREATE TABLE Category (
                              Name        varchar(128) DEFAULT 'Uncategorised' NOT NULL,
                              Description varchar(512),
                              PRIMARY KEY (Name));
    CREATE TABLE UserAddress (
                                 ID            varchar(64) NOT NULL,
                                 UserAccountID varchar(64) NOT NULL,
                                 PostalCode    int(4) DEFAULT 0 NOT NULL,
                                 Street        varchar(128) NOT NULL,
                                 PRIMARY KEY (ID,
                                              UserAccountID,
                                              PostalCode));
    CREATE TABLE PostalCode (
                                PostalCode int(4) DEFAULT 0 NOT NULL,
                                Area       varchar(128) DEFAULT 'Sted' NOT NULL,
                                PRIMARY KEY (PostalCode));
    CREATE TABLE Discount (
                              ID          varchar(64) NOT NULL,
                              Name        varchar(128) NOT NULL,
                              Description varchar(512),
                              Percentage  float NOT NULL,
                              Active      bit(1) NOT NULL,
                              PRIMARY KEY (ID));
    CREATE TABLE CartItem (
                              UserAccountID varchar(64) NOT NULL,
                              ProductID     varchar(64) NOT NULL,
                              Quantity      int(10) NOT NULL,
                              PRIMARY KEY (UserAccountID,
                                           ProductID),
                              INDEX (UserAccountID),
                              INDEX (ProductID));
    CREATE TABLE OrderItem (
                               ProductID      varchar(64) NOT NULL,
                               ProductOrderID varchar(64) NOT NULL,
                               Quantity       int(10) NOT NULL,
                               SubTotal       double NOT NULL,
                               PRIMARY KEY (ProductID,
                                            ProductOrderID),
                               INDEX (ProductID),
                               INDEX (ProductOrderID));
    CREATE TABLE ProductDiscount (
                                     DiscountID varchar(64) NOT NULL,
                                     ProductID  varchar(64) NOT NULL,
                                     PRIMARY KEY (DiscountID,
                                                  ProductID));
    ALTER TABLE ProductOrder ADD CONSTRAINT FKProductOrd309597 FOREIGN KEY (UserAccountID) REFERENCES UserAccount (ID) ON UPDATE Cascade;
    ALTER TABLE Product ADD CONSTRAINT FKProduct428805 FOREIGN KEY (BrandName) REFERENCES Brand (Name) ON UPDATE Cascade ON DELETE SET NULL ;
    ALTER TABLE Product ADD CONSTRAINT FKProduct125872 FOREIGN KEY (CategoryName) REFERENCES Category (Name) ON UPDATE Cascade ON DELETE SET DEFAULT;
    ALTER TABLE UserAddress ADD CONSTRAINT FKUserAddres739739 FOREIGN KEY (PostalCode) REFERENCES PostalCode (PostalCode) ON UPDATE Cascade ON DELETE No action;
    ALTER TABLE CartItem ADD CONSTRAINT FKCartItem790005 FOREIGN KEY (ProductID) REFERENCES Product (ID) ON UPDATE Cascade ON DELETE Cascade;
    ALTER TABLE OrderItem ADD CONSTRAINT FKOrderItem198854 FOREIGN KEY (ProductOrderID) REFERENCES ProductOrder (ID) ON UPDATE Cascade ON DELETE Cascade;
    ALTER TABLE OrderItem ADD CONSTRAINT FKOrderItem578114 FOREIGN KEY (ProductID) REFERENCES Product (ID) ON UPDATE Cascade ON DELETE Restrict;
    ALTER TABLE ProductDiscount ADD CONSTRAINT FKProductDis834859 FOREIGN KEY (DiscountID) REFERENCES Discount (ID) ON UPDATE Cascade ON DELETE Cascade;
    ALTER TABLE ProductDiscount ADD CONSTRAINT FKProductDis275846 FOREIGN KEY (ProductID) REFERENCES Product (ID) ON UPDATE Cascade ON DELETE Cascade;
    ALTER TABLE UserAddress ADD CONSTRAINT FKUserAddres929921 FOREIGN KEY (UserAccountID) REFERENCES UserAccount (ID) ON UPDATE Cascade ON DELETE Cascade;
    ALTER TABLE CartItem ADD CONSTRAINT FKCartItem480152 FOREIGN KEY (UserAccountID) REFERENCES UserAccount (ID) ON UPDATE Cascade ON DELETE Cascade;
    ALTER TABLE ProductOrder ADD CONSTRAINT FKProductOrd992331 FOREIGN KEY (PostalCode) REFERENCES PostalCode (PostalCode) ON UPDATE No action ON DELETE No action;
