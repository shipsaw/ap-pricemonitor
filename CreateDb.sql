CREATE TABLE IF NOT EXISTS Product (
                                       ProductID INT,
                                       Name TEXT NOT NULL,
                                       URL TEXT,
                                       Current_Price INT,
                                       Lowest_Price INT,
                                       Company INT NOT NULL DEFAULT 3,
                                       UNIQUE(URL)
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
