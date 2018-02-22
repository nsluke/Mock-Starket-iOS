//
//  LoginViewController.swift
//  Mock Starket
//
//  Created by Luke Solomon on 2/19/18.
//  Copyright © 2018 Luke Solomon. All rights reserved.
//

import Foundation
import Starscream
import Crashlytics

import Answers

class LoginViewController: ViewController {
    
    @IBOutlet weak var usernameTextField: UITextField!
    @IBOutlet weak var passwordTextField: UITextField!
    @IBOutlet weak var loginButton: UIButton!
    @IBOutlet weak var createAccountButton: UIButton!
    
    var socket = WebSocket(url: URL(string: "ws://159.89.154.221:8000/ws")!)

    
    override func viewDidLoad() {
        super.viewDidLoad()
        
        UIApplication.shared.statusBarStyle = .lightContent
        
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
        let message = "{\"action\": \"login\", \"value\": {\"username\": \"\(username)\", \"password\": \"\(password)\" }}"
        
        socket.disableSSLCertValidation = true
        print(message)
        socket.write(string: message)
    }
}

extension LoginViewController: WebSocketDelegate {
    func websocketDidConnect(socket: WebSocketClient) {
        print("Connected")
        loginButton.isUserInteractionEnabled = true
    }
    
    func websocketDidDisconnect(socket: WebSocketClient, error: Error?) {
        
    }
    
    func websocketDidReceiveMessage(socket: WebSocketClient, text: String) {
        print(text)
        
        let notificationView = UIView.init(frame: CGRect(x: 25, y: self.view.frame.height, width: self.loginButton.frame.width, height: self.loginButton.frame.height))
        notificationView.backgroundColor = UIColor.darkGray
        notificationView.cornerRadius = 5
        notificationView.shadowColor = UIColor.init(red: 1.0,
                                                    green: 1.0,
                                                    blue: 1.0,
                                                    alpha: 0.1)
        
        let rect = CGRect(x: notificationView.frame.midX,
                          y: notificationView.frame.midY,
                          width: self.loginButton.frame.width,
                          height: self.loginButton.frame.height)
        let notificationLabel = UILabel.init(frame: rect)
        notificationLabel.text = "Login Succesful!"
        notificationLabel.font = UIFont.init(name: "MuktaVaani", size: 20.0)
        notificationLabel.textColor = UIColor.white
        
        notificationView.addSubview(notificationLabel)
        self.view.addSubview(notificationView)
        
        UIView.animate(withDuration: 0.5, animations: {
            notificationView.frame = CGRect.init(x: 25, y: self.loginButton.frame.minY - 40, width: self.loginButton.frame.width, height: self.loginButton.frame.height)
            Answers.logCustomEvent(withName: "LogInSuccessful", customAttributes: ["any":"something"])

        }) { (bool) in
            sleep(1)
//            self.performSegue(withIdentifier: "loginSuccessful", sender: self)
        }
        
        
        
    }
    
    func websocketDidReceiveData(socket: WebSocketClient, data: Data) {
        print(data)
    }
}
