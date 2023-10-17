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
    "102t GLW Bogie Tanks" INT,
    "BAA/BZA Wagon Pack" INT,
    "BDA 80t Bogie Bolsters" INT,
    "BR Blue Diesel Electric Pack" INT,
    "BR Class 101" INT,
    "BR Class 150/1" INT,
    "BR Class 303" INT,
    "BR Class 31" INT,
    "BR Class 33" INT,
    "BR Class 421" INT,
    "BR Class 423" INT,
    "BR Class 73 'Gatwick Express'" INT,
    "BR Regional Railways Class 101" INT,
    "BR Sectors Class 56" INT,
    "Cargowaggon Flat IGA" INT,
    "Cargowaggon IWB" INT,
    "Chat Moss" INT,
    "Chatham Main & Medway Valley Lines" INT,
    "Chatham Main Line" INT,
    "Chatham Main Line: London Victoria & Blackfriars - Dover & Ramsgate" INT,
    "Class 150/1 Enhancement Pack" INT,
    "Class 150/2 Diesel Multiple Unit Pack" INT,
    "Class 153 Advanced" INT,
    "Class 156 Diesel Multiple Unit Pack" INT,
    "Class 158 (Perkins) Enhancement Pack" INT,
    "Class 158/159 (Cummins) Enhancement Pack" INT,
    "Class 168/170/171 Enhancement Pack" INT,
    "Class 170" INT,
    "Class 175 Enhancement Pack 2.0" INT,
    "Class 185 Multiple Unit Pack" INT,
    "Class 20 Advanced Collection" INT,
    "Class 205 Diesel Electric Multiple Unit Pack" INT,
    "Class 222 Advanced" INT,
    "Class 31 Enhancement Pack" INT,
    "Class 313 Electric Multiple Unit Pack" INT,
    "Class 314/315 Electric Multiple Unit Pack" INT,
    "Class 317 Electric Multiple Unit Pack Vol. 1" INT,
    "Class 317 Electric Multiple Unit Pack Vol. 2" INT,
    "Class 319 Electric Multiple Unit Pack Vol. 1" INT,
    "Class 319 Electric Multiple Unit Pack Vol. 2" INT,
    "Class 321 Electric Multiple Unit Pack" INT,
    "Class 325" INT,
    "Class 325 Enhancement Pack" INT,
    "Class 350 Enhancement Pack" INT,
    "Class 365 Enhancement Pack" INT,
    "Class 37 Locomotive Pack Vol. 1" INT,
    "Class 37 Locomotive Pack Vol. 2" INT,
    "Class 375/377 Enhancement Pack" INT,
    "Class 377/379/387 Enhancement Pack" INT,
    "Class 390" INT,
    "Class 390 Sound Pack" INT,
    "Class 411/412 Electric Multiple Unit Pack" INT,
    "Class 43 (MTU)/Mk3 Enhancement Pack" INT,
    "Class 43 (VP185)/Mk3 Enhancement Pack" INT,
    "Class 43 (Valenta)/Mk3 Enhancement Pack" INT,
    "Class 444/450 Enhancement Pack" INT,
    "Class 455 Enhancement Pack Vol. 1" INT,
    "Class 455 Enhancement Pack Vol. 2" INT,
    "Class 456 Electric Multiple Unit Pack" INT,
    "Class 465/466 Enhancement Pack Vol. 1" INT,
    "Class 465/466 Enhancement Pack Vol. 2" INT,
    "Class 50 Locomotive Pack" INT,
    "Class 56 Enhancement Pack" INT,
    "Class 60 Advanced" INT,
    "Class 66 Enhancement Pack" INT,
    "Class 67 Enhancement Pack" INT,
    "Class 68 Enhancement Pack" INT,
    "Class 700/707/717 Enhancement Pack" INT,
    "Class 800-803 Enhancement Pack" INT,
    "Class 86" INT,
    "Class 86 Enhancement Pack" INT,
    "Class 87 Locomotive Pack" INT,
    "Class 90 (Freightliner) Pack" INT,
    "Class 90/Mk3 DVT Pack" INT,
    "Class 91/Mk4 Enhancement Pack" INT,
    "Cloud Enhancement Pack" INT,
    "DB Schenker Class 59/2" INT,
    "ECML London to Peterborough" INT,
    "EWS & Freightliner Class 08" INT,
    "EWS Class 66 V2.0" INT,
    "EWS Class 67" INT,
    "EWS Class 92" INT,
    "European Community Asset Pack" INT,
    "FSA/FTA Wagon Pack" INT,
    "Freightliner Class 66 V2.0" INT,
    "Freightliner Class 70" INT,
    "GEML Class 90" INT,
    "GEML London to Ipswich" INT,
    "Gatwick Express Class 442" INT,
    "Grand Central Class 180" INT,
    "HEA Hoppers - Post BR" INT,
    "HHA Wagon Pack" INT,
    "HIA Wagon Pack" INT,
    "HKA/JMA Wagon Pack" INT,
    "ICA-D Wagon Pack" INT,
    "Intercity Class 91" INT,
    "Isle of Wight" INT,
    "JGA-K/PHA Wagon Pack" INT,
    "JHA Wagon Pack" INT,
    "JJA Autoballaster Advanced" INT,
    "JNA-C Wagon Pack" INT,
    "JPA Wagon Pack" INT,
    "JSA Wagon Pack" INT,
    "JTA/JUA/PTA Wagon Pack" INT,
    "JXA/POA Wagon Pack" INT,
    "Liverpool to Manchester" INT,
    "London Overground Class 378 - More Information OR DTG North London Line" INT,
    "London to Brighton" INT,
    "London to Faversham High Speed" INT,
    "MEA/PNA-F Wagon Pack" INT,
    "MFA/MHA/MTA Wagon Pack" INT,
    "MGR Wagon Pack" INT,
    "MKA/POA/ZKA Wagon Pack" INT,
    "Midland Main Line: London St. Pancras to Bedford" INT,
    "Midland Mainline: Derby - Leicester - Nottingham Extension" INT,
    "Midland Mainline: Sheffield - Derby" INT,
    "Mk1 Coach Pack Vol. 1" INT,
    "Mk2A-C Coach Pack" INT,
    "Mk2D-F Coach Pack" INT,
    "Mk2F DBSO Coach Pack" INT,
    "Mk3A-B Coach Pack" INT,
    "Network SouthEast Class 159" INT,
    "North London & Goblin Lines" INT,
    "North London Line - More Information OR DTG North London & Goblin Lines" INT,
    "North Wales Coast Line: Crewe to Holyhead" INT,
    "Portsmouth Direct Line" INT,
    "Portsmouth Direct Line: London Waterloo - Portsmouth" INT,
    "Powerhaul Class 66 V2.0" INT,
    "Riviera Line: Exeter to Paignton" INT,
    "ScotRail Class 68" INT,
    "Settle to Carlisle" INT,
    "Signal Enhancement Pack" INT,
    "Sky & Weather Enhancement Pack 2.0" INT,
    "South London Network" INT,
    "South Wales Coastal: Bristol to Swansea" INT,
    "South Western Expressways - Reading" INT,
    "South Western Main Line: Southampton to Bournemouth" INT,
    "Southeastern Class 465" INT,
    "Southern Class 455/8" INT,
    "TDA-D Wagon Pack" INT,
    "TIA-Y Wagon Pack" INT,
    "TTA Wagon Pack Vol. 1" INT,
    "Thameslink Class 700" INT,
    "Track Enhancement Pack" INT,
    "VGA/VKA Wagon" INT,
    "Voyager Advanced" INT,
    "WCML South: London Euston - Birmingham" INT,
    "Wagon (4 Wheel) Sound Pack" INT,
    "Wagon (Flat) Sound Pack" INT,
    "Wagon (Modern) Sound Pack" INT,
    "Wagon (Old) Sound Pack" INT,
    "West Coast Main Line North" INT,
    "West Coast Main Line Over Shap" INT,
    "Western Lines Scotland" INT,
    "Wherry Lines: Norwich to Great Yarmouth & Lowestoft Route 2.0" INT,
    "Wherry Lines: Norwich to Great Yarmouth & Lowestoft Route 2.0 - UPGRADE" INT,
    "YGB Seacow Advanced" INT,
    "YQA Parr" INT,
    "ZCA Sea Urchins" INT
)
