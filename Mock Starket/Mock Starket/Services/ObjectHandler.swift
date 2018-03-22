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
    
    public var stockArray = [Stock]()
    
    var stockSet:NSMutableOrderedSet = NSMutableOrderedSet() //.addObjects(from: ObjectHandler.sharedInstance.stockArray) as! NSMutableOrderedSet
    
    public var marketArray = [Stock]()
    var marketSet:NSMutableOrderedSet = NSMutableOrderedSet()
    
    
    // User info
    var currentUserName:String?
    var currentUserID:String?
    var netWorth = Double()
    var displayName = String()
    var wallet = Double()
    
    
    // ======================== init ======================== //
    override init() {
        super.init()
    }
    
    init(currentUser: User) {
        super.init()
        
//        for stock in stockArray {
//            stockSet.add(stock.name)
//        }
    }
    
    // ======================== Action ======================== //
    func actionRouter (json: JSON) {
        /* actions can either be a single message or an array, so we must check which we have by using an optional bind. If it's an array, send the singleActionRout an array of actions. If not, send the rout only one. */
        if let actionArray = json.array {
            for action in actionArray {
                singleActionRout(action:action)
            }
        } else {
            singleActionRout(action:json)
        }
    }
    
    func singleActionRout(action:JSON) {
        let actionType = action["action"].stringValue
        if actionType == "object" {
            objectHandler(message: action["msg"])
        } else if actionType == "update" {
            /* Messages on an update can either be a single message or an array, so we must check which we have by using an optional bind. If it's an array, send the update handler an array of update messages. If not, send the update handler only one. */
            if let messageArray = action["msg"].array {
                for message in messageArray {
                    updateHandler(message:message)
                }
            } else {
                updateHandler(message:action["msg"])
            }
        } else if actionType == "alert" {
            alertHandler(message:action["msg"])
        } else if actionType == "chat" {
            chatHandler(message:action["msg"])
        } else if actionType == "login" {
            loginHandler(message: action["msg"])
        }
    }
    
    // ======================== Object ======================== //
    func objectHandler (message: JSON) {
        let objectType = message["type"].stringValue
        let objectID = message["id"].stringValue
        let object = message["object"]
        
        if objectType == "user" {
            if objectID == self.currentUserID! {
                self.displayName = object["display_name"].stringValue
            } else {
                //Todo: other users
            }
            
        } else if objectType == "stock" {
            let stock = Stock.init(name: objectID, value: object["current_price"].doubleValue, fullname: object["name"].stringValue)
            marketSet.add(stock.name)
            marketArray.append(stock)
            
            NotificationCenter.default.post(name: ObjectServiceNotification.ObjectStock.rawValue,
                                            object: nil,
                                            userInfo: ["id": objectID, "value": object["current_price"].doubleValue])

        } else if objectType == "portfolio" {
            if objectID == self.currentUserID! {
                self.wallet = object["wallet"].doubleValue
                self.netWorth = object["net_worth"].doubleValue
            } else {
                //TODO: other users
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
            
            // ========================== //
            // Update user
            // ========================== //
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
            // ========================== //
            // Update portfolio
            // ========================== //
            } else if updateType == "portfolio" {
                
                // ========================== //
                // Update to the current user
                // ========================== //
                if updateID == self.currentUserID! {
                    if changeType == "wallet" {
                        changeValue = change["value"].doubleValue
                        self.wallet = changeValue as! Double
                        
                        NotificationCenter.default.post(name: ObjectServiceNotification.ActionUpdateCurrentUserPortfolioWallet.rawValue,
                                                        object: nil,
                                                        userInfo: ["id": updateID, "value": changeValue])
                    } else if changeType == "net_worth" {
                        changeValue = change["value"].doubleValue
                        self.netWorth = changeValue as! Double
                        
                        NotificationCenter.default.post(name: ObjectServiceNotification.ActionUpdateCurrentUserPortfolioNetWorth.rawValue,
                                                        object: nil,
                                                        userInfo: ["id": updateID, "value": changeValue])
                    } else if changeType == "ledger" {
                        //TODO:
                        //                        changeValue = change["value"]
                        //
                        //                        NotificationCenter.default.post(name: ObjectServiceNotification.ActionUpdateCurrentUserPortfolioLedger.rawValue,
                        //                                                        object: nil,
                        //                                                        userInfo: ["id": updateID, "value": changeValue])
                    }
                
                // ========================== //
                // Update to other user
                // ========================== //
                } else {
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
                }
            // ========================== //
            // Update stock
            // ========================== //
            } else if updateType == "stock" {
                // ========================== //
                // Update to the market
                // ========================== //
                changeValue = change["value"].doubleValue
                
                let stock = Stock.init(name: updateID, value: changeValue as! Double)
                let index = marketSet.index(of: updateID)
                var amountChanged = 0.0
                
                if changeValue as! Double != 0 && index <= marketArray.count - 1 && marketArray[index].value != 0 {
                    amountChanged = changeValue as! Double - marketArray[index].value
                }
                
                if self.marketSet.contains(updateID) {
                    marketArray.remove(at: index)
                    marketArray.insert(stock, at: index)
                    if marketArray[index].recordValue < stock.value {
                        marketArray[index].recordValue = stock.value
                    }
                    marketArray[index].amountChanged = amountChanged
                } else {
                    marketSet.add(stock.name)
                    marketArray.append(stock)
                }
                
                NotificationCenter.default.post(name: ObjectServiceNotification.ActionUpdateStockPrice.rawValue,
                                                object: nil,
                                                userInfo: ["id": updateID, "value": changeValue])
                
                // ========================== //
                // Update to the current user
                // ========================== //
                let userIndex = stockSet.index(of: updateID)
                var userAmountChanged = 0.0
                
                if changeValue as! Double != 0 && index <= stockArray.count - 1 && stockArray[index].value != 0 {
                    amountChanged = changeValue as! Double - stockArray[index].value
                }
                
                if self.stockSet.contains(updateID) {
                    stockArray.remove(at: index)
                    stockArray.insert(stock, at: index)
                    if stockArray[index].recordValue < stock.value {
                        stockArray[index].recordValue = stock.value
                    }
                    
                    stockArray[index].amountChanged = amountChanged
                    
                    NotificationCenter.default.post(name: ObjectServiceNotification.ActionUpdateCurrentUserStockPrice.rawValue,
                                                    object: nil,
                                                    userInfo: ["id": updateID, "value": changeValue])
                }
                

                

            // ========================== //
            // Update ledger
            // ========================== //
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
        let alertDict = ["alert" : Alert.init(json:message)]
        NotificationCenter.default.post(name: ObjectServiceNotification.ActionChat.rawValue,
                                        object: nil,
                                        userInfo: alertDict)
    }
    
    // ======================== Login ======================== //
    func loginHandler (message: JSON) {
        if message["success"].boolValue {
            self.currentUserID = message["uuid"].stringValue
            NotificationCenter.default.post(name: ObjectServiceNotification.LoginSuccessful.rawValue,
                                            object: nil,
                                            userInfo: nil)
        } else {
            
            let errorDict = ["err" : message["err"].stringValue]

            NotificationCenter.default.post(name: ObjectServiceNotification.LoginUnsuccessful.rawValue,
                                            object: nil,
                                            userInfo: errorDict)
        }
    }
}

enum ObjectServiceNotification: Notification.Name {
    //Object updates
    case ObjectStock = "ObjectStock"
    
    //user updates
    case ActionUpdateUsername = "ActionUpdateUsername"
    case ActionUpdateUserActive = "ActionUpdateUserActive"
    //currentUser updates
    case ActionUpdateCurrentUserPortfolioWallet = "ActionUpdateCurrentUserPortfolioWallet"
    case ActionUpdateCurrentUserPortfolioNetWorth = "ActionUpdateCurrentUserPortfolioNetWorth"
    case ActionUpdateCurrentUserPortfolioLedger = "ActionUpdateCurrentUserPortfolioLedger"
    //portfolio updates
    case ActionUpdatePortfolioWallet = "ActionUpdatePortfolioWallet"
    case ActionUpdatePortfolioNetWorth = "ActionUpdatePortfolioNetWorth"
    case ActionUpdatePortfolioLedger = "ActionUpdatePortfolioLedger"
    //stock updates
    case ActionUpdateStockPrice = "ActionUpdateStockPrice"
    case ActionUpdateCurrentUserStockPrice = "ActionUpdateCurrentUserStockPrice"
    //exchange updates
    case ActionUpdateExchangeHolders = "ActionUpdateExchangeHolders"
    case ActionUpdateExchangeOpenShares = "ActionUpdateExchangeOpenShares"
    
    //chat updates
    case ActionChat = "ActionChat"
    
    //alert updates
    case ActionAlert = "ActionAlert"
    
    //login
    case LoginSuccessful = "LoginSuccesful"
    case LoginUnsuccessful = "LoginUnsuccesful"
}
