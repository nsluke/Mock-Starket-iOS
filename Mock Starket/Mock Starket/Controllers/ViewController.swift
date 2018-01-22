//
//  ViewController.swift
//  Mock Starket
//
//  Created by Luke Solomon on 1/15/18.
//  Copyright © 2018 Luke Solomon. All rights reserved.
//

import UIKit
import RxSwift
import RxCocoa
import RxStarscream
import Starscream

class ViewController: UIViewController {

    @IBOutlet weak var textField: UITextField!
    @IBOutlet weak var textview: UITextView!
    
    @IBOutlet weak var sendButton: UIButton!
    
    private let disposeBag = DisposeBag()
    private let socket:WebSocket = WebSocket(url: URL(string: "ws://159.203.244.103:8000/ws")!)
    private let writeSubject = PublishSubject<String>()
    
    
    override func viewDidLoad() {
        super.viewDidLoad()
        
        let responseString = socket.rx.response.map { response -> String in
            
            print(response)
            
            switch response {
                case .connected:
                    print(response)
                    return "Connected\n"
                
                case .disconnected(let error):
                    return "Disconnected with error: \(String(describing: error)) \n"
                case .message(let msg):
                    return "RESPONSE (Message): \(msg) \n"
                case .data(let data):
                    return  "RESPONSE (Data): \(data) \n"
                case .pong:
                    return "RESPONSE (Pong)"
            }
        }
        
        Observable.merge([responseString, writeSubject.asObservable()])
            .scan([]) { lastMsg, newMsg -> Array<String> in
                return Array(lastMsg + [newMsg])
            }.map { $0.joined(separator: "\n")
            }.asDriver(onErrorJustReturn: "")
            .drive(textview.rx.text)
            .disposed(by: disposeBag)
        socket.disableSSLCertValidation = true
        
        socket.connect()
        //self.sendMessage(message:"{\"action\": \"login\", \"value\": {\"username\": \"username\", \"password\":\"password\"}}" )
        self.textField.text = "{\"action\": \"login\", \"value\": {\"username\": \"username\", \"password\":\"password\"}}"
    }


    @objc fileprivate func sendMessage(message: String) {
        socket.write(string: message)
        writeSubject.onNext("SENT: \(message) \n")
        textField.text = nil
        textField.resignFirstResponder()
    }
    
    @IBAction func sendButtonTapped(_ sender: Any) {
        if let text = textField.text {
            sendMessage(message:text)

        }
    }
    
    
    
}

//extension ViewController: UITextFieldDelegate {
//
//    func textFieldShouldReturn(_ textField: UITextField) -> Bool {
////        if let text = textField.text, !text.isEmpty {
////            sendMessage(message: text)
////            return true
////        }
//        return false
//    }
//}

