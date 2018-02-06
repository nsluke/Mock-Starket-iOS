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
        self.value = value.init(json["action"]["value"])
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
    var ticker: [String : [String:Int]]
    
    init(json: JSON) {
        self.ticker = json.arrayValue
    }
}
