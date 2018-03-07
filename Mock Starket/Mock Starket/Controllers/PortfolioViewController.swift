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
    
    //MARK: IBOutlets
    @IBOutlet weak var headerView: UIView!
    @IBOutlet weak var sideMenuButton: UIButton!
    @IBOutlet weak var blueView: UIView!
    @IBOutlet weak var tableViewHeaderView: UIView!
    @IBOutlet weak var portfolioStampView: UIView!
    @IBOutlet weak var tableView: UITableView!
    
    //header outlets
    @IBOutlet weak var netWorthLabel: UILabel!
    @IBOutlet weak var netWorthPercentageChangeLabel: UILabel!
    @IBOutlet weak var netWorthPercentSignLabel: UILabel!
    @IBOutlet weak var netWorthArrowIcon: UIImageView!
    @IBOutlet weak var netWorthDollarSignLabel: UILabel!
    
    
    //portfolio stamp files
    @IBOutlet weak var cashAmountLabel: UILabel!
    @IBOutlet weak var investmentsAmountLabel: UILabel!
    
    
//    case CHUNT = "Chunt's Hats"
//    case KING = "Paddle King"
//    case CBIO = "Sebio's Streaming Services"
//    case OW = "Overwatch"
//    case SCOTT = "Michael Scott Paper Company"
    
//    case DM = "Dunder Milf"
//    case GWEN = "Gwent"
//    case CHU = "Chu Supply"
//    case SWEET = "Sweet Sweet Tea"
//    case TRAP = "❤ Trap 4 Life"
    
//    case FIG = "Figgis Agency"
//    case ZONE = "Danger Zone"
//    case PLNX = "Planet Express"
//    case MOM = "Mom's Friendly Robot Company"
    
    var portfolioArray = [
        Stock.init(name: "CHUNT", value: 0.0),
        Stock.init(name: "KING", value: 0.0),
        Stock.init(name: "CBIO", value: 0.0),
        Stock.init(name: "OW", value: 0.0),
        Stock.init(name: "SCOTT", value: 0.0),
        
        Stock.init(name: "DM", value: 0.0),
        Stock.init(name: "GWEN", value: 0.0),
        Stock.init(name: "CHU", value: 0.0),
        Stock.init(name: "SWEET", value: 0.0),
        Stock.init(name: "TRAP", value: 0.0),
        
        Stock.init(name: "FIG", value: 0.0),
        Stock.init(name: "ZONE", value: 0.0),
        Stock.init(name: "PLNX", value: 0.0),
        Stock.init(name: "MOM", value: 0.0)
        ]
    var mutableSet = NSMutableOrderedSet()
    var netWorth = Double()
    
    
    //MARK: View Lifecycle
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
    
    override func viewDidDisappear(_ animated: Bool) {
        NotificationCenter.default.removeObserver(self, name: NetworkServiceNotification.SocketMessageReceived.rawValue, object: nil)
    }
    
    //MARK: View Setup Functions
    func setupViews() {
        UIApplication.shared.statusBarStyle = .lightContent

        let gradient = CAGradientLayer.init()
        gradient.colors = [UIColor.init(red: 1.0/20.0, green: 1.0/30.0, blue: 1.0/48.0, alpha: 0.0),
                           UIColor.init(red: 1.0/20.0, green: 1.0/30.0, blue: 1.0/48.0, alpha: 1.0)]
        gradient.startPoint = CGPoint.init(x: blueView.frame.width/2, y: blueView.frame.minY)
        gradient.endPoint = CGPoint.init(x: blueView.frame.width/2, y: blueView.frame.maxY)
        blueView.layer.addSublayer(gradient)
        
        self.netWorthLabel.text = "Loading..."
        self.netWorthPercentageChangeLabel.text = ""
        self.netWorthPercentSignLabel.text = ""
        self.netWorthArrowIcon.isHidden = true
        self.netWorthDollarSignLabel.isHidden = true
        
        self.cashAmountLabel.text = "Loading..."
        self.investmentsAmountLabel.text = "Loading..."
        
        for i in portfolioArray {
            mutableSet.add(i.name)
        }
        
        tableView.reloadData()
    }
    
    //MARK: NotificationHandling
    @objc func update(_ notification:NSNotification) {
        guard let actionArray = notification.userInfo?["actionArray"] as? [ResponseAction] else {
            return
        }
        
        // Array for initial objects
        // Set for knowing if the string is there < if we don't use a stock we can just use ordered
        // Ordered Set for knowing the index
        
        for action in actionArray {
            if action.action == "update" && action.type == "stock"{
                for change in action.changes {
                    if change.field == "current_price" {
                        let stock = Stock.init(name: action.id, value: change.value)
                        let index = mutableSet.index(of: stock.name)
                        var amountChanged = 0.0
                        
                        if change.value != 0 && portfolioArray[index].value != 0 {
                             amountChanged = change.value - portfolioArray[index].value
                        }
                        
                        if mutableSet.contains(stock.name) {
                            portfolioArray.remove(at: index)
                            portfolioArray.insert(stock, at: index)
                            if portfolioArray[index].recordValue < stock.value {
                                portfolioArray[index].recordValue = stock.value
                            }
                            
                            portfolioArray[index].amountChanged = amountChanged
                        } else {
                            mutableSet.add(stock.name)
                            portfolioArray.append(stock)
                            
                            
                        }
                        
                        self.tableView.reloadData()
                    } else {
                        print("New Field!" + change.field)
                    }
                }
            } else if action.action == "update" && action.type == "portfolio" && action.id == "1" {
                
                for change in action.changes {
                    if change.field == "net_worth" {
                        //Handle net worth change
                        if self.netWorth == 0 {
                            self.netWorth = change.value
                            self.netWorthDollarSignLabel.isHidden = false
                        }
                        
                        let percentChange = round((((change.value - self.netWorth) / self.netWorth) * 100 ) * 100) / 100
                        print(percentChange)
                        
                        if percentChange > 0 {
                            self.netWorthLabel.text = String(format: "%.2f", change.value)
                            self.netWorthPercentageChangeLabel.text = String(format: "%.2f", percentChange)
                            self.netWorthPercentSignLabel.text = "%"
                            
                            self.netWorthLabel.textColor = UIColor.msAquamarine
                            self.netWorthPercentSignLabel.textColor = UIColor.msAquamarine
                            self.netWorthPercentageChangeLabel.textColor = UIColor.msAquamarine
                            self.netWorthArrowIcon.isHidden = false
                            
                            self.netWorthArrowIcon.image = UIImage.init(imageLiteralResourceName: "uptriangle")
                            
                            
                        } else if percentChange == 0 {
                            self.netWorthLabel.text = String(format: "%.2f", change.value)
                            self.netWorthPercentageChangeLabel.text = ""
                            self.netWorthPercentSignLabel.text = ""
                            
                            self.netWorthLabel.textColor = UIColor.white
                            self.netWorthPercentSignLabel.textColor = UIColor.white
                            self.netWorthPercentageChangeLabel.textColor = UIColor.white
                            
                            self.netWorthArrowIcon.isHidden = true
                        } else if percentChange < 0 {
                            self.netWorthLabel.text = String(format: "%.2f", change.value)
                            self.netWorthPercentageChangeLabel.text = String(format: "%.2f", percentChange)
                            self.netWorthPercentSignLabel.text = "%"
                            
                            self.netWorthLabel.textColor = UIColor.msFlatRed
                            self.netWorthPercentSignLabel.textColor = UIColor.msFlatRed
                            self.netWorthPercentageChangeLabel.textColor = UIColor.msFlatRed
                            
                            self.netWorthArrowIcon.isHidden = false
                            self.netWorthArrowIcon.image = UIImage.init(imageLiteralResourceName: "downtriangle")
                        }
                        self.netWorth = change.value
                    }
                }
            
            }
        }
    }
    
    //Mark: IBActions
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
        cell?.nameLabel.text = portfolioArray[indexPath.row].fullname

        if portfolioArray[indexPath.row].value <= 0 {
            cell?.costLabel.text = ""
        } else {
            cell?.costLabel.text = String(format: "%.2f", portfolioArray[indexPath.row].value)
        }
        
        if portfolioArray[indexPath.row].value <= 0 {
            cell?.recordLabel.text = ""
        } else {
            cell?.recordLabel.text = String(format: "%.2f", portfolioArray[indexPath.row].recordValue)
        }
        
        let amountChanged = portfolioArray[indexPath.row].amountChanged
        cell?.changeLabel.text = String(format: "%.2f", amountChanged)

        if amountChanged < 0 {
            cell?.changeLabel.textColor = UIColor.msFlatRed
            cell?.changeImageView.isHidden = false
            cell?.changeImageView.image = #imageLiteral(resourceName: "downtriangle")
        } else if amountChanged == 0.0 {
            cell?.changeLabel.text = ""
            cell?.changeLabel.textColor = UIColor.msLightGray
            cell?.changeImageView.isHidden = true
        } else if amountChanged > 0 {
            cell?.changeLabel.textColor = UIColor.msAquamarine
            cell?.changeImageView.isHidden = false
            cell?.changeImageView.image = #imageLiteral(resourceName: "uptriangle")
        }
        
        
        
        return cell!
    }
}


