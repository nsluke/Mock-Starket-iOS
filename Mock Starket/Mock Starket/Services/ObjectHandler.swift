//
//  ObjectHandler.swift
//  Mock Starket
//
//  Created by Luke Solomon on 3/15/18.
//  Copyright © 2018 Luke Solomon. All rights reserved.
//

import Foundation
import SwiftyJSON

final class ObjectHandler: NSObject {
    public static let sharedInstance = ObjectHandler()
    
    public var stockArray = [
        Stock.init(name: "CHUNT", value: 0.0),
        Stock.init(name: "KING", value: 0.0),
        Stock.init(name: "CBIO", value: 0.0),
        Stock.init(name: "OW", value: 0.0),
        Stock.init(name: "SCOTT", value: 0.0),

        Stock.init(name: "DM", value: 0.0),
        Stock.init(name: "GWEN", value: 0.0),
        Stock.init(name: "CHU", value: 0.0),
        Stock.init(name: "SWEET", value: 0.0),
        Stock.init(name: "TRAP", value: 0.0),

        Stock.init(name: "FIG", value: 0.0),
        Stock.init(name: "ZONE", value: 0.0),
        Stock.init(name: "PLNX", value: 0.0),
        Stock.init(name: "MOM", value: 0.0)
    ]
    var stockSet = NSMutableOrderedSet()
    
    // User info
    var currentUser:User
    var netWorth = Double()

    
    // ======================== init ======================== //
    override init() {
        for stock in stockArray {
            stockSet.add(stock.name)
        }
        
        
        
    }
    

    // ======================== Action ======================== //
    func actionRouter (json: JSON) {
        //actions come in as an array, so we parse through them to see what kind of action we're dealing with: Update, Alert, or Chat
        for action in json.arrayValue {
            let actionType = action["action"].stringValue
            if actionType == "update" {
                /* Messages on an update can either be a single message or an array, so we must check which we have by using an optional bind. If it's an array, send the update handler an array of update messages. If not, send the update handler only one. */
                if let actionArray = action["msg"].array {
                    for message in actionArray {
                        updateHandler(message:message)
                    }
                } else {
                    updateHandler(message:action["msg"])
                }
            } else if actionType == "alert" {
                alertHandler(message:json["msg"])
            } else if actionType == "chat" {
                chatHandler(message:json["msg"])
            }
        }
    }
    
    
    // ======================== Update ======================== //
    func updateHandler (message: JSON) {
        //Update come in one of four types, User, Portfolio, Stock, and Exchange, so we have to determine which we're receiving, and post the appropriate notification with that information.
        let updateType = message["type"].stringValue
        let updateID = message["id"].stringValue
        
        for change in message["changes"].arrayValue {
            let changeType = change["field"].stringValue
            let changeValue:Any
            
            if updateType == "user" {
                if changeType == "display_name" {
                    changeValue = change["value"].stringValue
                    NotificationCenter.default.post(name: ObjectServiceNotification.ActionUpdateUsername.rawValue,
                                                    object: nil,
                                                    userInfo: ["id": updateID, "value": changeValue])
                } else if changeType == "active" {
                    changeValue = change["value"].boolValue
                    
                    NotificationCenter.default.post(name: ObjectServiceNotification.ActionUpdateUserActive.rawValue,
                                                    object: nil,
                                                    userInfo: ["id": updateID, "value": changeValue])
                }
            } else if updateType == "portfolio" {
                if changeType == "wallet" {
                    changeValue = change["value"].doubleValue
                    NotificationCenter.default.post(name: ObjectServiceNotification.ActionUpdatePortfolioWallet.rawValue,
                                                    object: nil,
                                                    userInfo: ["id": updateID, "value": changeValue])
                } else if changeType == "net_worth" {
                    changeValue = change["value"].doubleValue
                    
                    NotificationCenter.default.post(name: ObjectServiceNotification.ActionUpdatePortfolioNetWorth.rawValue,
                                                    object: nil,
                                                    userInfo: ["id": updateID, "value": changeValue])
                } else if changeType == "ledger" {
                    //TODO:
//                    changeValue = change["value"]
//
//                    NotificationCenter.default.post(name: ObjectServiceNotification.ActionUpdatePortfolioLedger.rawValue,
//                                                    object: nil,
//                                                    userInfo: ["id": updateID, "value": changeValue])
                }
            } else if updateType == "stock" {
                changeValue = change["current_price"].doubleValue
                
                let stock = Stock.init(name: updateID, value: changeValue as! Double)
                let index = stockSet.index(of: updateID)
                var amountChanged = 0.0
                
                if changeValue as! Double != 0 && stockArray[index].value != 0 {
                    amountChanged = changeValue as! Double - stockArray[index].value
                }
                
                if self.stockSet.contains(updateID) {
                    stockArray.remove(at: index)
                    stockArray.insert(stock, at: index)
                    if stockArray[index].recordValue < stock.value {
                        stockArray[index].recordValue = stock.value
                    }
                    stockArray[index].amountChanged = amountChanged
                } else {
                    stockSet.add(stock.name)
                    stockArray.append(stock)
                }
                
                NotificationCenter.default.post(name: ObjectServiceNotification.ActionUpdateStockPrice.rawValue,
                                                object: nil,
                                                userInfo: ["id": updateID, "value": changeValue])
                
            } else if updateType == "exchange_ledger" {
                if changeType == "holders" {
                    //TODO:
//                    changeValue = change["value"].doubleValue
//                    NotificationCenter.default.post(name: ObjectServiceNotification.ActionUpdateExchangeHolders.rawValue,
//                                                    object: nil,
//                                                    userInfo: ["id": updateID, "value": changeValue])
                } else if changeType == "open_shares" {
                    changeValue = change["value"].intValue
                    NotificationCenter.default.post(name: ObjectServiceNotification.ActionUpdateExchangeOpenShares.rawValue,
                                                    object: nil,
                                                    userInfo: ["id": updateID, "value": changeValue])
                }
            } else {
                print("unexpected update type!")
            }
        }
    }
    
    
    // ======================== Chat ======================== //
    func chatHandler (message: JSON) {
        let chatDict = ["chat" : Chat.init(json: message)]
        NotificationCenter.default.post(name: ObjectServiceNotification.ActionChat.rawValue,
                                        object: nil,
                                        userInfo: chatDict)
    }
    
    
    // ======================== Alert ======================== //
    func alertHandler (message: JSON) {
        let alertDict = ["chat" : Alert.init(json:message)]
        NotificationCenter.default.post(name: ObjectServiceNotification.ActionChat.rawValue,
                                        object: nil,
                                        userInfo: alertDict)
    }
}

enum ObjectServiceNotification: Notification.Name {
    //user updates
    case ActionUpdateUsername = "ActionUpdateUsername"
    case ActionUpdateUserActive = "ActionUpdateUserActive"
    //portfolio updates
    case ActionUpdatePortfolioWallet = "ActionUpdatePortfolioWallet"
    case ActionUpdatePortfolioNetWorth = "ActionUpdatePortfolioNetWorth"
    case ActionUpdatePortfolioLedger = "ActionUpdatePortfolioLedger"
    //stock updates
    case ActionUpdateStockPrice = "ActionUpdateStockPrice"
    //exchange updates
    case ActionUpdateExchangeHolders = "ActionUpdateExchangeHolders"
    case ActionUpdateExchangeOpenShares = "ActionUpdateExchangeOpenShares"

    //chat updates
    case ActionChat = "ActionChat"
    
    //alert updates
    case ActionAlert = "ActionAlert"
}
