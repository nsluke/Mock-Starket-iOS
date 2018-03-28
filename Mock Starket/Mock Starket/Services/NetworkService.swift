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
    
//    public static let socket = WebSocket(url: URL(string: "ws://localhost:8000/ws")!)
    public static let socket = WebSocket(url: URL(string: "ws://159.89.154.221:8000/ws")!)
    
    private override init() {
        super.init()
    }
    
    static func connect() {
        socket.delegate = sharedInstance
        socket.disableSSLCertValidation = true
        socket.connect()
    }
    
    static func login(username: String, password: String) {
        let message = "{\"action\": \"login\", \"msg\": {\"username\": \"\(username)\", \"password\": \"\(password)\" }}"
        print(message)
        
        socket.write(string: message)
    }
    
    static func createAccount(username:String, password: String) {

    
        let message = "{\"action\": \"new_account\", \"msg\": {\"username\": \"\(username)\", \"password\": \"\(password)\" }}"
        print(message)
        
        socket.write(string: message)
        
    }
    
    static func buyStock(ticker:String, amount:Int) {
        
        let message = "{\"action\": \"trade\", \"msg\": {\"stock_ticker\": \"\(ticker)\", \"amount\": \(amount) }}"
        print(message)
        
        socket.write(string: message)
        
    }
    
    static func sellStock(ticker:String, amount:Int) {
        let message = "{\"action\": \"trade\", \"msg\": {\"stock_ticker\": \"\(ticker)\", \"amount\": -\(amount) }}"
        print(message)
        
        socket.write(string: message)
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
        
        // Parse the socket response from a string into JSON, then send it
        // to the ObjectHandler to be routed
        
        let json = JSON.init(parseJSON: text)
//        dump(json)
        print(json)
        
        NotificationCenter.default.post(name: NetworkServiceNotification.SocketMessageReceived.rawValue,
                                        object: nil,
                                        userInfo: nil)
        
        ObjectHandler.sharedInstance.actionRouter(json: json)
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
