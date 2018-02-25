//
//  PortfolioViewController.swift
//  Mock Starket
//
//  Created by Luke Solomon on 2/20/18.
//  Copyright © 2018 Luke Solomon. All rights reserved.
//

import UIKit
import SideMenu
import Starscream
import Crashlytics
import SwiftyJSON


class PortfolioViewController: UIViewController {
    
    @IBOutlet weak var headerView: UIView!
    @IBOutlet weak var sideMenuButton: UIButton!
    @IBOutlet weak var blueView: UIView!
    @IBOutlet weak var tableViewHeaderView: UIView!
    @IBOutlet weak var portfolioStampView: UIView!
    @IBOutlet weak var tableView: UITableView!
    
    var socket = WebSocket(url: URL(string: "ws://159.89.154.221:8000/ws")!)
    
    
    override func viewDidLoad() {
        super.viewDidLoad()
        
        let menuRightNavigationController = UISideMenuNavigationController(rootViewController: self)
        SideMenuManager.default.menuRightNavigationController = menuRightNavigationController
        SideMenuManager.default.menuAddPanGestureToPresent(toView: sideMenuButton)
//        SideMenuManager.default.menuAddScreenEdgePanGesturesToPresent(toView: self.navigationController!.view)
        SideMenuManager.default.menuAddScreenEdgePanGesturesToPresent(toView: self.navigationController!.view, forMenu: UIRectEdge.all)
    }
    
    override func viewDidAppear(_ animated: Bool) {
        super.viewDidAppear(true)
        self.setupViews()
    }
    
    func setupViews() {
        UIApplication.shared.statusBarStyle = .lightContent

        let gradient = CAGradientLayer.init()
        gradient.colors = [UIColor.init(red: 1.0/20.0, green: 1.0/30.0, blue: 1.0/48.0, alpha: 0.0), 
                           UIColor.init(red: 1.0/20.0, green: 1.0/30.0, blue: 1.0/48.0, alpha: 1.0) ]
        gradient.startPoint = CGPoint.init(x: blueView.frame.width/2, y: blueView.frame.minY)
        gradient.endPoint = CGPoint.init(x: blueView.frame.width/2, y: blueView.frame.minY)
        blueView.layer.addSublayer(gradient)
    }
    
    //handle button tap
    @IBAction func sideMenuButtonTapped(_ sender: Any) {
        self.present(SideMenuManager.default.menuRightNavigationController!, animated: true, completion: nil)
        dismiss(animated: true, completion: nil)
    }
    
}

extension PortfolioViewController: UITabBarDelegate {
    func tabBar(_ tabBar: UITabBar, didSelect item: UITabBarItem) {
        
    }
}


extension PortfolioViewController: UITableViewDataSource {
    func numberOfSections(in tableView: UITableView) -> Int {
        return 1
    }
    
    func tableView(_ tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
        return 1
    }
    
    func tableView(_ tableView: UITableView, cellForRowAt indexPath: IndexPath) -> UITableViewCell {
        let cell = tableView.dequeueReusableCell(withIdentifier: "portfolioTableViewCell", for: indexPath)
        cell.textLabel?.text = "Congrats, this is all the app does for now. Thanks for helping!"
        
        return cell
    }
}

extension PortfolioViewController: WebSocketDelegate {
    func websocketDidConnect(socket: WebSocketClient) {
        
    }
    
    func websocketDidDisconnect(socket: WebSocketClient, error: Error?) {
        
    }
    
    func websocketDidReceiveMessage(socket: WebSocketClient, text: String) {
        
        let action = Action.init(json: JSON.parse(text))
        print(action)
        
    }
    
    func websocketDidReceiveData(socket: WebSocketClient, data: Data) {
        print(data)

    }
}
