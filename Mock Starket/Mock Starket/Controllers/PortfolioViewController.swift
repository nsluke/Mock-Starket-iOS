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
    
    var portfolioArray = [Stock]()
    var mutableSet = NSMutableOrderedSet()
    
    
    override func viewDidLoad() {
        super.viewDidLoad()
        
//        let menuLeftNavigationController = storyboard!.instantiateViewController(withIdentifier: "LeftMenuNavigationController") as! UISideMenuNavigationController
//        let menuRightNavigationController = storyboard!.instantiateViewController(withIdentifier: "RightMenuNavigationController") as! UISideMenuNavigationController
//        SideMenuManager.default.menuLeftNavigationController = menuLeftNavigationController
//        SideMenuManager.default.menuRightNavigationController = menuRightNavigationController
//        SideMenuManager.default.menuAddPanGestureToPresent(toView: sideMenuButton)
//        SideMenuManager.default.menuAddScreenEdgePanGesturesToPresent(toView: self.navigationController!.view, forMenu: UIRectEdge.all)

        NotificationCenter.default.addObserver(self, selector: #selector(update(_:)), name: NetworkServiceNotification.SocketMessageReceived.rawValue, object: nil)
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
        gradient.endPoint = CGPoint.init(x: blueView.frame.width/2, y: blueView.frame.maxY)
        blueView.layer.addSublayer(gradient)
    }
    
    @objc func update(_ notification:NSNotification) {
        guard let response = notification.userInfo?["actionArray"] as? [ResponseAction] else {
            return
        }
        
        // Array for initial objects
        // Set for knowing if the string is there < if we don't use a stock we can just use ordered
        // Ordered Set for knowing the index
        
        
        for action in response {
            if action.action == "update" && action.type == "stock"{
                for change in  action.changes {
                    if change.field == "current_price" {
                        let stock = Stock.init(name: action.id, value: change.value)
                        
                        if mutableSet.contains(stock.name) {
                            let index = mutableSet.index(of: stock.name)
                            portfolioArray.remove(at: index)
                            portfolioArray.insert(stock, at: index)
                        } else {
                            mutableSet.add(stock.name)
                            portfolioArray.append(stock)
                        }

                        self.tableView.reloadData()
                    } else {
                        print("New Field!" + change.field)
                    }
                }
            }
        }
        
    }
    
    //handle button tap
    @IBAction func sideMenuButtonTapped(_ sender: Any) {
        self.present(SideMenuManager.default.menuRightNavigationController!, animated: true, completion: nil)
        dismiss(animated: true, completion: nil)
    }
    
}

extension PortfolioViewController: UITabBarDelegate {
    func tabBar(_ tabBar: UITabBar, didSelect item: UITabBarItem) {
        // Code for cool animation goes here
    }
}

extension PortfolioViewController: UITableViewDataSource {
    func tableView(_ tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
        return self.mutableSet.count
    }
    
    func tableView(_ tableView: UITableView, cellForRowAt indexPath: IndexPath) -> UITableViewCell {
        let cell = tableView.dequeueReusableCell(withIdentifier: "portfolioTableViewCell", for: indexPath) as? PortfolioTableViewCell
        
        cell?.tickerLabel.text = portfolioArray[indexPath.row].name
        cell?.costLabel.text = String(format: "%.2f", portfolioArray[indexPath.row].value)
        cell?.recordLabel.text = ""
        cell?.changeLabel.text = ""
        cell?.nameLabel.text = ""
        
        return cell!
    }
}


