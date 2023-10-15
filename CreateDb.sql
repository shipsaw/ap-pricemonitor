CREATE TABLE IF NOT EXISTS Product (
                                       ProductID INT,
                                       Name TEXT NOT NULL UNIQUE,
                                       URL TEXT COLLATE NOCASE,
                                       Current_Price INT,
                                       Lowest_Price INT,
                                       Company INT NOT NULL DEFAULT 5
    );

CREATE TABLE IF NOT EXISTS EssentialJoin (
                                             ProductID INT,
                                             EssentialID INT,
                                             FOREIGN KEY(ProductID) REFERENCES Product(ROWID),
    FOREIGN KEY(EssentialID) REFERENCES Product(ROWID)
    );

CREATE TABLE IF NOT EXISTS ScenarioJoin (
                                            ProductID INT,
                                            ScenarioID INT,
                                            FOREIGN KEY(ProductID) REFERENCES Product(ROWID),
    FOREIGN KEY(ScenarioID) REFERENCES Product(ROWID)
    );
