CREATE TABLE IF NOT EXISTS NewProduct (
    ProductID INT,
    Name TEXT NOT NULL UNIQUE,
    URL TEXT COLLATE NOCASE,
    Current_Price INT,
    Lowest_Price INT,
    Company INT NOT NULL DEFAULT 5
);

CREATE TABLE IF NOT EXISTS NewEssentialJoin (
    ProductID INT,
    EssentialID INT,
    FOREIGN KEY(ProductID) REFERENCES Product(ROWID),
    FOREIGN KEY(EssentialID) REFERENCES Product(ROWID)
);

CREATE TABLE IF NOT EXISTS NewScenarioJoin (
    ProductID INT,
    ScenarioID INT,
    FOREIGN KEY(ProductID) REFERENCES Product(ROWID),
    FOREIGN KEY(ScenarioID) REFERENCES Product(ROWID)
);

CREATE TABLE IF NOT EXISTS NewRecommendedJoin (
    ProductID INT,
    ScenarioID INT,
    FOREIGN KEY(ProductID) REFERENCES Product(ROWID),
    FOREIGN KEY(ScenarioID) REFERENCES Product(ROWID)
);

CREATE TABLE IF NOT EXISTS PriceReporting(
    Date TEXT UNIQUE NOT NULL,
    "102t GLW Bogie Tanks",
    "BAA/BZA Wagon Pack",
    "BDA 80t Bogie Bolsters",
    "BR Blue Diesel Electric Pack",
    "BR Class 101",
    "BR Class 150/1",
    "BR Class 303",
    "BR Class 31",
    "BR Class 33",
    "BR Class 421",
    "BR Class 423",
    "BR Class 73 'Gatwick Express'",
    "BR Regional Railways Class 101",
    "BR Sectors Class 56",
    "Cargowaggon Flat IGA",
    "Cargowaggon IWB",
    "Chat Moss",
    "Chatham Main & Medway Valley Lines",
    "Chatham Main Line",
    "Chatham Main Line: London Victoria & Blackfriars - Dover & Ramsgate",
    "Class 150/1 Enhancement Pack",
    "Class 150/2 Diesel Multiple Unit Pack",
    "Class 153 Advanced",
    "Class 156 Diesel Multiple Unit Pack",
    "Class 158 (Perkins) Enhancement Pack",
    "Class 158/159 (Cummins) Enhancement Pack",
    "Class 168/170/171 Enhancement Pack",
    "Class 170",
    "Class 175 Enhancement Pack 2.0",
    "Class 185 Multiple Unit Pack",
    "Class 20 Advanced Collection",
    "Class 205 Diesel Electric Multiple Unit Pack",
    "Class 222 Advanced",
    "Class 31 Enhancement Pack",
    "Class 313 Electric Multiple Unit Pack",
    "Class 314/315 Electric Multiple Unit Pack",
    "Class 317 Electric Multiple Unit Pack Vol. 1",
    "Class 317 Electric Multiple Unit Pack Vol. 2",
    "Class 319 Electric Multiple Unit Pack Vol. 1",
    "Class 319 Electric Multiple Unit Pack Vol. 2",
    "Class 321 Electric Multiple Unit Pack",
    "Class 325",
    "Class 325 Enhancement Pack",
    "Class 350 Enhancement Pack",
    "Class 365 Enhancement Pack",
    "Class 37 Locomotive Pack Vol. 1",
    "Class 37 Locomotive Pack Vol. 2",
    "Class 375/377 Enhancement Pack",
    "Class 377/379/387 Enhancement Pack",
    "Class 390",
    "Class 390 Sound Pack",
    "Class 411/412 Electric Multiple Unit Pack",
    "Class 43 (MTU)/Mk3 Enhancement Pack",
    "Class 43 (VP185)/Mk3 Enhancement Pack",
    "Class 43 (Valenta)/Mk3 Enhancement Pack",
    "Class 444/450 Enhancement Pack",
    "Class 455 Enhancement Pack Vol. 1",
    "Class 455 Enhancement Pack Vol. 2",
    "Class 456 Electric Multiple Unit Pack",
    "Class 465/466 Enhancement Pack Vol. 1",
    "Class 465/466 Enhancement Pack Vol. 2",
    "Class 50 Locomotive Pack",
    "Class 56 Enhancement Pack",
    "Class 60 Advanced",
    "Class 66 Enhancement Pack",
    "Class 67 Enhancement Pack",
    "Class 68 Enhancement Pack",
    "Class 700/707/717 Enhancement Pack",
    "Class 800-803 Enhancement Pack",
    "Class 86",
    "Class 86 Enhancement Pack",
    "Class 87 Locomotive Pack",
    "Class 90 (Freightliner) Pack",
    "Class 90/Mk3 DVT Pack",
    "Class 91/Mk4 Enhancement Pack",
    "Cloud Enhancement Pack",
    "DB Schenker Class 59/2",
    "ECML London to Peterborough",
    "EWS & Freightliner Class 08",
    "EWS Class 66 V2.0",
    "EWS Class 67",
    "EWS Class 92",
    "European Community Asset Pack",
    "FSA/FTA Wagon Pack",
    "Freightliner Class 66 V2.0",
    "Freightliner Class 70",
    "GEML Class 90",
    "GEML London to Ipswich",
    "Gatwick Express Class 442",
    "Grand Central Class 180",
    "HEA Hoppers - Post BR",
    "HHA Wagon Pack",
    "HIA Wagon Pack",
    "HKA/JMA Wagon Pack",
    "ICA-D Wagon Pack",
    "Intercity Class 91",
    "Isle of Wight",
    "JGA-K/PHA Wagon Pack",
    "JHA Wagon Pack",
    "JJA Autoballaster Advanced",
    "JNA-C Wagon Pack",
    "JPA Wagon Pack",
    "JSA Wagon Pack",
    "JTA/JUA/PTA Wagon Pack",
    "JXA/POA Wagon Pack",
    "Liverpool to Manchester",
    "London Overground Class 378 - More Information OR DTG North London Line",
    "London to Brighton",
    "London to Faversham High Speed",
    "MEA/PNA-F Wagon Pack",
    "MFA/MHA/MTA Wagon Pack",
    "MGR Wagon Pack",
    "MKA/POA/ZKA Wagon Pack",
    "Midland Main Line: London St. Pancras to Bedford",
    "Midland Mainline: Derby - Leicester - Nottingham Extension",
    "Midland Mainline: Sheffield - Derby",
    "Mk1 Coach Pack Vol. 1",
    "Mk2A-C Coach Pack",
    "Mk2D-F Coach Pack",
    "Mk2F DBSO Coach Pack",
    "Mk3A-B Coach Pack",
    "Network SouthEast Class 159",
    "North London & Goblin Lines",
    "North London Line - More Information OR DTG North London & Goblin Lines",
    "North Wales Coast Line: Crewe to Holyhead",
    "Portsmouth Direct Line",
    "Portsmouth Direct Line: London Waterloo - Portsmouth",
    "Powerhaul Class 66 V2.0",
    "Riviera Line: Exeter to Paignton",
    "ScotRail Class 68",
    "Settle to Carlisle",
    "Signal Enhancement Pack",
    "Sky & Weather Enhancement Pack 2.0",
    "South London Network",
    "South Wales Coastal: Bristol to Swansea",
    "South Western Expressways - Reading",
    "South Western Main Line: Southampton to Bournemouth",
    "Southeastern Class 465",
    "Southern Class 455/8",
    "TDA-D Wagon Pack",
    "TIA-Y Wagon Pack",
    "TTA Wagon Pack Vol. 1",
    "Thameslink Class 700",
    "Track Enhancement Pack",
    "VGA/VKA Wagon",
    "Voyager Advanced",
    "WCML South: London Euston - Birmingham",
    "Wagon (4 Wheel) Sound Pack",
    "Wagon (Flat) Sound Pack",
    "Wagon (Modern) Sound Pack",
    "Wagon (Old) Sound Pack",
    "West Coast Main Line North",
    "West Coast Main Line Over Shap",
    "Western Lines Scotland",
    "Wherry Lines: Norwich to Great Yarmouth & Lowestoft Route 2.0",
    "Wherry Lines: Norwich to Great Yarmouth & Lowestoft Route 2.0 - UPGRADE",
    "YGB Seacow Advanced",
    "YQA Parr",
    "ZCA Sea Urchins"
)
