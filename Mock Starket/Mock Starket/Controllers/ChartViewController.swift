//
//  ChartViewController.swift
//  Mock Starket
//
//  Created by Luke Solomon on 1/21/18.
//  Copyright © 2018 Luke Solomon. All rights reserved.
//

import UIKit
import RxSwift
import RxCocoa
import RxStarscream
import Starscream


class ChartViewController: UIViewController {
    
    private let disposeBag = DisposeBag()
    private var socket:WebSocket = WebSocket(url: URL(string: "ws://159.203.244.103:8000/ws")!)
    private let writeSubject = PublishSubject<String>()

    
    
    override func viewDidLoad() {
        super.viewDidLoad()
        
        socket = WebSocket(url: URL(string: "ws://localhost:8080/")!)
        socket.delegate = self as! WebSocketDelegate
        socket.connect()
        
        Observable.just([])


        // Do any additional setup after loading the view.
    }
    

    override func didReceiveMemoryWarning() {
        super.didReceiveMemoryWarning()
        // Dispose of any resources that can be recreated.
    }
    
    
}
