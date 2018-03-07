//
//  NetworkService.swift
//  Mock Starket
//
//  Created by Luke Solomon on 1/15/18.
//  Copyright © 2018 Luke Solomon. All rights reserved.
//

import Foundation
import Starscream
import SwiftyJSON


final class NetworkService: NSObject {
    public static let sharedInstance = NetworkService()
    
    public static let socket = WebSocket(url: URL(string: "ws://localhost:8000/ws")!)

//    public static let socket = WebSocket(url: URL(string: "ws://159.89.154.221:8000/ws")!)
    
    private override init() {
        super.init()
    }
    
    static func connect() {
        socket.delegate = sharedInstance
        socket.disableSSLCertValidation = true
        socket.connect()
    }
    
    static func login(username: String, password: String) {
        let message = "{\"action\": \"login\", \"value\": {\"username\": \"\(username)\", \"password\": \"\(password)\" }}"
        print(message)
        
        socket.write(string: message)
    }

    static func receiveSocket () {
        
    }
}
extension NetworkService: WebSocketDelegate {
    func websocketDidConnect(socket: WebSocketClient) {
        print("Connected")
        
        NotificationCenter.default.post(name: NetworkServiceNotification.SocketDidConnect.rawValue,
                                        object: nil,
                                        userInfo: nil)
    }
    
    func websocketDidDisconnect(socket: WebSocketClient, error: Error?) {
        print("Disconnected")
        
        NotificationCenter.default.post(name: NetworkServiceNotification.SocketDidDisconnect.rawValue,
                                        object: nil,
                                        userInfo: nil)
    }
    
    func websocketDidReceiveMessage(socket: WebSocketClient, text: String) {
        debugPrint(text, separator: "%n")
        
        let json = JSON.init(parseJSON: text)
        var actionArray = [Action]()
        
        for action in json.arrayValue { // run through the actionArray
            actionArray.append(Action.init(json: action))
        }
        
        debugPrint("actionArray: \(actionArray)", separator: "%n")
        
        NotificationCenter.default.post(name: NetworkServiceNotification.SocketMessageReceived.rawValue,
                                        object: text,
                                        userInfo: ["actionArray" : actionArray])
    }
    
    func websocketDidReceiveData(socket: WebSocketClient, data: Data) {
        print(data)
    }
}


enum NetworkServiceNotification: Notification.Name {
    case SocketMessageReceived = "SocketMessageReceived"
    case SocketDidConnect = "SocketDidConnect"
    case SocketDidDisconnect = "SocketDidDisconnect"
}

extension Notification.Name: ExpressibleByStringLiteral {
    public init(stringLiteral value: String) {
        self.init(value)
    }
    
    public init(extendedGraphemeClusterLiteral value: String) {
        self.init(value)
    }
    
    public init(unicodeScalarLiteral value: String) {
        self.init(value)
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
