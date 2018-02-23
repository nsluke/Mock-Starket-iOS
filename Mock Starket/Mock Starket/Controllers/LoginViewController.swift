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

class LoginViewController: ViewController {
    
    //MARK: IBOutlets
    @IBOutlet weak var scrollView: UIScrollView!
    @IBOutlet weak var usernameTextField: UITextField!
    @IBOutlet weak var passwordTextField: UITextField!
    @IBOutlet weak var loginButton: UIButton!
    @IBOutlet weak var createAccountButton: UIButton!
    @IBOutlet weak var loginActivityIndicator: UIActivityIndicatorView!
    @IBOutlet weak var constraintContentHeight: NSLayoutConstraint!

    //MARK: Variables
    var activeTextField:UITextField!
    var socket = WebSocket(url: URL(string: "ws://159.89.154.221:8000/ws")!)
    var lastOffset: CGPoint!
    var keyboardHeight: CGFloat!
    
    
    //MARK: View Lifecycle
    override func viewDidLoad() {
        super.viewDidLoad()
        
        UIApplication.shared.statusBarStyle = .lightContent
        loginButton.backgroundColor = UIColor.red
        self.loginActivityIndicator.isHidden = true
        NotificationCenter.default.addObserver(self, selector: #selector(keyboardWillShow(notification:)), name: NSNotification.Name.UIKeyboardWillShow, object: nil)
        NotificationCenter.default.addObserver(self, selector: #selector(keyboardWillHide(notification:)), name: NSNotification.Name.UIKeyboardWillHide, object: nil)
        
        socket.delegate = self
        socket.disableSSLCertValidation = true
        socket.connect()
    }
    
    //MARK: IBActions
    @IBAction func touchReceived(_ sender: Any) {
        guard activeTextField != nil else { return }
        
        activeTextField?.resignFirstResponder()
        activeTextField = nil
    }
    
    @IBAction func loginButtonTapped(_ sender: Any) {
        guard let username = usernameTextField.text else { return }
        guard let password = passwordTextField.text else { return }
        
        let message = "{\"action\": \"login\", \"value\": {\"username\": \"\(username)\", \"password\": \"\(password)\" }}"
        print(message)
        self.loginButton.titleLabel?.attributedText = NSAttributedString(string: "")
        self.loginActivityIndicator.isHidden = false
        self.loginActivityIndicator.startAnimating()
        
        socket.write(string: message)
    }
    
    //MARK: keyboard scroll code
    @objc func keyboardWillShow(notification: NSNotification) {
        if keyboardHeight != nil {
            return
        }
        if let keyboardSize = (notification.userInfo?[UIKeyboardFrameBeginUserInfoKey] as? NSValue)?.cgRectValue {
            keyboardHeight = keyboardSize.height
            // so increase contentView's height by keyboard height
            UIView.animate(withDuration: 0.3, animations: {
                self.constraintContentHeight.constant += self.keyboardHeight
            })
            // move if keyboard hide input field
            let distanceToBottom = self.scrollView.frame.size.height - (activeTextField?.frame.origin.y)! - (activeTextField?.frame.size.height)!
            let collapseSpace = keyboardHeight - distanceToBottom
            if collapseSpace < 0 {
                // no collapse
                return
            }
            // set new offset for scroll view
            UIView.animate(withDuration: 0.3, animations: {
                // scroll to the position above keyboard 10 points
                self.scrollView.contentOffset = CGPoint(x: self.lastOffset.x, y: collapseSpace + 10)
            })
        }
    }
    
    @objc func keyboardWillHide(notification: NSNotification) {
        UIView.animate(withDuration: 0.3) {
            self.constraintContentHeight.constant -= self.keyboardHeight
            self.scrollView.contentOffset = self.lastOffset
        }
        keyboardHeight = nil
    }
    
    //MARK: textField Checks
    func checkForErrors() -> Bool {
        var errors = false
        let title = "Error"
        var message = ""
        
        if usernameTextField.text!.isEmpty {
            errors = true
            message += "First name empty"
            alertWithTitle(title, message: message, ViewController: self, toFocus:self.usernameTextField)
        } else if passwordTextField.text!.isEmpty {
            errors = true
            message += "Surname empty"
            alertWithTitle(title, message: message, ViewController: self, toFocus:self.passwordTextField)
            
        } else if !isValidEmail(usernameTextField.text!) {
            errors = true
            message += "Invalid Email Address"
            alertWithTitle(title, message: message, ViewController: self, toFocus:self.usernameTextField)
        } else if (passwordTextField.text?.count)! < 8 {
            errors = true
            message += "Password must be at least 8 characters"
            alertWithTitle(title, message: message, ViewController: self, toFocus:self.passwordTextField)
        }
        return errors
    }
    
    func isValidEmail(_ test:String) -> Bool {
        // email validation here...
        return true
    }
    
    func alertWithTitle(_ title: String!, message: String, ViewController: UIViewController, handler: @escaping (UIAlertAction) -> (Void)) {
        let alert = UIAlertController(title: title, message: message, preferredStyle: .alert)
        let action = UIAlertAction(title: "OK", style: UIAlertActionStyle.cancel,handler: handler);
        alert.addAction(action)
        ViewController.present(alert, animated: true, completion:nil)
    }
    
    func alertWithTitle(_ title: String!, message: String, ViewController: UIViewController,toFocus:UITextField) {
        let alert = UIAlertController(title: title, message: message, preferredStyle: .alert)
        let action = UIAlertAction(title: "OK", style: UIAlertActionStyle.cancel,handler: {_ in
            toFocus.becomeFirstResponder()
        });
        alert.addAction(action)
        ViewController.present(alert, animated: true, completion:nil)
    }
    
}

extension LoginViewController: UITextFieldDelegate {
    func textFieldShouldBeginEditing(_ textField: UITextField) -> Bool {
        self.activeTextField = textField
        lastOffset = self.scrollView.contentOffset
        return true
    }
    
    func textFieldShouldReturn(_ textField: UITextField) -> Bool {
        activeTextField.resignFirstResponder()
        activeTextField = nil
        
        if (textField == self.usernameTextField) {
            self.passwordTextField.becomeFirstResponder()
        } else {
            let thereWereErrors = checkForErrors()
            if !thereWereErrors {
                loginButton.backgroundColor = UIColor.red
            } else {
                loginButton.backgroundColor = UIColor(red: 26.0/255.0, green: 188.0/255.0, blue: 156.0/255.0, alpha: 1.0)
            }
        }
        
        return true
    }
}

extension LoginViewController: WebSocketDelegate {
    func websocketDidConnect(socket: WebSocketClient) {
        print("Connected")
        loginButton.backgroundColor = UIColor(red: 26.0/255.0, green: 188.0/255.0, blue: 156.0/255.0, alpha: 1.0)
    }
    
    func websocketDidDisconnect(socket: WebSocketClient, error: Error?) {
        
    }
    
    func websocketDidReceiveMessage(socket: WebSocketClient, text: String) {
        print(text)
        
        self.loginActivityIndicator.stopAnimating()
        self.loginActivityIndicator.isHidden = true
        
        Answers.logCustomEvent(withName: "LogInSuccessful", customAttributes: ["any":"something"])
        alertWithTitle("Login Successful!", message: "", ViewController: self) { (alertAction) -> (Void) in
            self.performSegue(withIdentifier: "loginSuccessful", sender: self)
        }
    }
    
    func websocketDidReceiveData(socket: WebSocketClient, data: Data) {
        print(data)
    }
}
