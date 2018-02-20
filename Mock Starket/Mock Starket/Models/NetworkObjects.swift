//
//  NetworkObjects.swift
//  Mock Starket
//
//  Created by Luke Solomon on 1/27/18.
//  Copyright © 2018 Luke Solomon. All rights reserved.
//

import Foundation
import SwiftyJSON


struct Action {
    var type:String
    var value:Value
    
    init(json:JSON) {
        self.type = json["action"].stringValue
        self.value = Value.init(json: json["action"]["value"])
    }
}

struct Value {
    var type:String
    var object:Any
    
    init(json: JSON) {
        self.type = json["type"].stringValue
        
        if self.type == "portfolio" {
            self.object = Portfolio.init(json: json["value"])
        } else if self.type == "valuable" {
            self.object = Valuable.init(json: json["value"])
        } else {
            object = 0
        } 
    }
}

struct Portfolio {
    var name:String
    var uuid:Int
    var wallet:Float
    var net_worth:Float
    var ledger:Ledger
    
    init(json:JSON) {
        self.name = json["object"]["name"].stringValue
        self.uuid = json["object"]["uuid"].intValue
        self.wallet = json["object"]["net_worth"].floatValue
        self.net_worth = 0
        self.ledger = Ledger.init(json: json["object"]["ledger"])
    }
}

struct Valuable {
    var name:String
    var tickerID:String
    var current_price:Float
    
    init(json:JSON) {
        self.name = json["name"].stringValue
        self.tickerID = json["ticker_id"].stringValue
        self.current_price = json["current_price"].floatValue
    }
}

struct Ledger {
    var ticker: Stock
    
    init(json: JSON) {
//        for i in json {
        self.ticker = Stock.init(name: json["name"].stringValue, value: json["amount"].doubleValue) // this part is gonna break
//        }
    }
}

struct Stock {
    var name: String
    var value: Double
    
    init(name:String, value:Double) {
        self.name = name
        self.value = value
    }
    
}
