//
//  LoginViewController.swift
//  Mock Starket
//
//  Created by Luke Solomon on 2/19/18.
//  Copyright © 2018 Luke Solomon. All rights reserved.
//

import Foundation
import Starscream

class LoginViewController: ViewController {
    
    @IBOutlet weak var usernameTextField: UITextField!
    @IBOutlet weak var passwordTextField: UITextField!
    @IBOutlet weak var loginButton: UIButton!
    @IBOutlet weak var createAccountButton: UIButton!
    
    var socket = WebSocket(url: URL(string: "ws://159.89.154.221:8000/ws")!)

    
    override func viewDidLoad() {
        super.viewDidLoad()
        
        loginButton.isUserInteractionEnabled = false
        
        socket.delegate = self
        socket.connect()
    }
    
    @IBAction func loginButtonTapped(_ sender: Any) {
        
        guard let username = usernameTextField.text else {
            return
        }
        
        guard let password = passwordTextField.text else {
            return
        }
        
        sendLoginRequest(username: username, password: password)
    }
    
    func sendLoginRequest(username:String, password:String) {
        let message = "{\"action\": \"login\", \"value\": {\"username\": \(username), \"password\": \(password)}}"
        
        socket.disableSSLCertValidation = true
        socket.write(string: message)
    }
    
}

extension LoginViewController: WebSocketDelegate {
    
    func websocketDidConnect(socket: WebSocketClient) {
        print("Connected")
        
        loginButton.isUserInteractionEnabled = true

//        socket.write(string: "{\"action\": \"login\", \"value\": {\"username\": \"username\", \"password\":\"password\"}}")
    }
    
    func websocketDidDisconnect(socket: WebSocketClient, error: Error?) {
        
    }
    
    func websocketDidReceiveMessage(socket: WebSocketClient, text: String) {
        print(text)
    }
    
    func websocketDidReceiveData(socket: WebSocketClient, data: Data) {
        print(data)
    }
    
}
