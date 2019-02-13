//
//  NetworkObjects.swift
//  Mock Starket
//
//  Created by Luke Solomon on 1/27/18.
//  Copyright © 2018 Luke Solomon. All rights reserved.
//

import Foundation
import SwiftyJSON

// ======================== Actions ======================== //
struct Update {
    var type:String
    var id:String
    var changes:[Change]
    
    init(json: JSON) {
        self.id = json["id"].stringValue
        self.type = json["type"].stringValue
        
        var changeArray = [Change]()
        for change in json["changes"].arrayValue {
            changeArray.append(Change.init(json:change))
        }
        self.changes = changeArray
    }
    
    init(id:String, type:String, changes:[Change]) {
        self.id = id
        self.type = type
        self.changes = changes
    }
}

struct Chat {
    var messageBody:String
    var author:String
    var timestamp:Int
    
    init(json: JSON) {
        self.messageBody = json["message_body"].stringValue
        self.author = json["author"].stringValue
        self.timestamp = json["timestamp"].intValue
    }
    
    init(messageBody: String, author: String, timestamp:Int) {
        self.messageBody = messageBody
        self.author = author
        self.timestamp = timestamp
    }
}

struct Alert {
    var alert:String
    var type:String
    var timestamp:Int

    init(json: JSON) {
        self.alert = json["alert"].stringValue
        self.type = json["type"].stringValue
        self.timestamp = json["timestamp"].intValue
    }
    
    init(alert: String, type: String, timestamp:Int) {
        self.alert = alert
        self.type = type
        self.timestamp = timestamp
    }
}

// ======================== Messages ======================== //
struct Message {
    var changes:[Change]
    
    init(json:JSON) {
        self.changes = [Change]()
        for change in json["msg"]["changes"].arrayValue {
            self.changes.append(Change.init(json: change))
        }
    }
}

// ======================== Changes ======================== //
struct Change {
    var field:String
    var value:Any
    
    init(json:JSON) {
        self.field = json["field"].stringValue
        
        if self.field == "ledger" {
            
        }
        
        self.value = json["value"].doubleValue
    }
}

// ======================== Models ======================== //
struct Exchange {
    var name:String
    var ledger: [String : LedgerEntry]
}

struct LedgerEntry {
    var portfolioID:String
    var stockID:String
    var amount:Int
    
    init(portfolioID:String, stockID:String, amount:Int) {
        self.portfolioID = portfolioID
        self.stockID = stockID
        self.amount = amount
    }
}

struct User {
    var id:String
    var displayName:String
    
    init (json:JSON) {
        self.id = json["id"].stringValue
        self.displayName = json["displayName"].stringValue
    }
    
    init (id:String, displayName:String) {
        self.id = id
        self.displayName = displayName
    }
    
}

struct Portfolio {
    var name:String
    var uuid:Int
    var wallet:Float
    var net_worth:Float
    
    init(_ json:JSON) {
        self.name = json["object"]["name"].stringValue
        self.uuid = json["object"]["uuid"].intValue
        self.wallet = json["object"]["net_worth"].floatValue
        self.net_worth = 0
    }
}

struct Stock {
    var name: String
    var fullname: String
//    var uuid: String
    var value: Double
    var recordValue: Double
    var amountChanged: Double
    
    init(name:String, value:Double) {
        self.name = name
        
        switch name {
        case "CHUNT":
            self.fullname = StockNames.CHUNT.rawValue
        case "KING":
            self.fullname = StockNames.KING.rawValue
        case "CBIO":
            self.fullname = StockNames.CBIO.rawValue
        case "OW":
            self.fullname = StockNames.OW.rawValue
        case "SCOTT":
            self.fullname = StockNames.SCOTT.rawValue
        case "DM":
            self.fullname = StockNames.DM.rawValue
        case "GWEN":
            self.fullname = StockNames.GWEN.rawValue
        case "CHU":
            self.fullname = StockNames.CHU.rawValue
        case "SWEET":
            self.fullname = StockNames.SWEET.rawValue
        case "TRAP":
            self.fullname = StockNames.TRAP.rawValue
        case "FIG":
            self.fullname = StockNames.FIG.rawValue
        case "ZONE":
            self.fullname = StockNames.ZONE.rawValue
        case "PLNX":
            self.fullname = StockNames.PLNX.rawValue
        case "MOM":
            self.fullname = StockNames.MOM.rawValue
        default:
            self.fullname = ""
        }
        
        self.value = value
        self.recordValue = value
        self.amountChanged = 0.0
    }
    
    init(name:String, value:Double, fullname:String) {
        self.name = name
        self.fullname = fullname
        self.value = value
        self.recordValue = value
        self.amountChanged = 0.0
    }
}

enum StockNames: String {
    case CHUNT = "Chunt's Hats"
    case KING = "Paddle King"
    case CBIO = "Sebio's Streaming Services"
    case OW = "Overwatch"
    case SCOTT = "Michael Scott Paper Company"
    case DM = "Dunder Milf"
    case GWEN = "Gwent"
    case CHU = "Chu Supply"
    case SWEET = "Sweet Sweet Tea"
    case TRAP = "❤ Trap 4 Life"
    case FIG = "Figgis Agency"
    case ZONE = "Danger Zone"
    case PLNX = "Planet Express"
    case MOM = "Mom's Friendly Robot Company"
}

/*
guard let actionArray = notification.userInfo?["actionArray"] as? [ResponseAction] else {
    return
}

// Array for initial objects
// Set for knowing if the string is there < if we don't use a stock we can just use ordered
// Ordered Set for knowing the index

for action in actionArray {
    if action.action == "update" && action.type == "stock"{
        for change in action.changes {
            if change.field == "current_price" {
                let stock = Stock.init(name: action.id, value: change.value)
                let index = mutableSet.index(of: stock.name)
                var amountChanged = 0.0
 
                if change.value != 0 && portfolioArray[index].value != 0 {
                    amountChanged = change.value - portfolioArray[index].value
                }
 
                if mutableSet.contains(stock.name) {
                    portfolioArray.remove(at: index)
                    portfolioArray.insert(stock, at: index)
                    if portfolioArray[index].recordValue < stock.value {
                        portfolioArray[index].recordValue = stock.value
                    }
 
                    portfolioArray[index].amountChanged = amountChanged
                } else {
                    mutableSet.add(stock.name)
                    portfolioArray.append(stock)
 
 
                }
 
                self.tableView.reloadData()
            } else {
                print("New Field!" + change.field)
            }
        }
    } else if action.action == "update" && action.type == "portfolio" && action.id == "1" {
 
        for change in action.changes {
            if change.field == "net_worth" {
                //Handle net worth change
                if self.netWorth == 0 {
                    self.netWorth = change.value
                    self.netWorthDollarSignLabel.isHidden = false
                }
 
                let percentChange = round((((change.value - self.netWorth) / self.netWorth) * 100 ) * 100) / 100
                print(percentChange)
 
                if percentChange > 0 {
                    self.netWorthLabel.text = String(format: "%.2f", change.value)
                    self.netWorthPercentageChangeLabel.text = String(format: "%.2f", percentChange)
                    self.netWorthPercentSignLabel.text = "%"
 
                    self.netWorthLabel.textColor = UIColor.msAquamarine
                    self.netWorthPercentSignLabel.textColor = UIColor.msAquamarine
                    self.netWorthPercentageChangeLabel.textColor = UIColor.msAquamarine
                    self.netWorthArrowIcon.isHidden = false
 
                    self.netWorthArrowIcon.image = UIImage.init(imageLiteralResourceName: "uptriangle")
 
 
                } else if percentChange == 0 {
                    self.netWorthLabel.text = String(format: "%.2f", change.value)
                    self.netWorthPercentageChangeLabel.text = ""
                    self.netWorthPercentSignLabel.text = ""
 
                    self.netWorthLabel.textColor = UIColor.white
                    self.netWorthPercentSignLabel.textColor = UIColor.white
                    self.netWorthPercentageChangeLabel.textColor = UIColor.white
 
                    self.netWorthArrowIcon.isHidden = true
                } else if percentChange < 0 {
                    self.netWorthLabel.text = String(format: "%.2f", change.value)
                    self.netWorthPercentageChangeLabel.text = String(format: "%.2f", percentChange)
                    self.netWorthPercentSignLabel.text = "%"
 
                    self.netWorthLabel.textColor = UIColor.msFlatRed
                    self.netWorthPercentSignLabel.textColor = UIColor.msFlatRed
                    self.netWorthPercentageChangeLabel.textColor = UIColor.msFlatRed
 
                    self.netWorthArrowIcon.isHidden = false
                    self.netWorthArrowIcon.image = UIImage.init(imageLiteralResourceName: "downtriangle")
                }
                self.netWorth = change.value
            }
        }
 
    }
}
*/
