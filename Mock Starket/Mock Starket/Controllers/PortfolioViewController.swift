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
    
    
    //MARK: View Lifecycle
    override func viewDidLoad() {
        super.viewDidLoad()
        
//        let menuLeftNavigationController = storyboard!.instantiateViewController(withIdentifier: "LeftMenuNavigationController") as! UISideMenuNavigationController
//        let menuRightNavigationController = storyboard!.instantiateViewController(withIdentifier: "RightMenuNavigationController") as! UISideMenuNavigationController
//        SideMenuManager.default.menuLeftNavigationController = menuLeftNavigationController
//        SideMenuManager.default.menuRightNavigationController = menuRightNavigationController
//        SideMenuManager.default.menuAddPanGestureToPresent(toView: sideMenuButton)
//        SideMenuManager.default.menuAddScreenEdgePanGesturesToPresent(toView: self.navigationController!.view, forMenu: UIRectEdge.all)
        
        NotificationCenter.default.addObserver(self,
                                               selector: #selector(updateWallet(_:)),
                                               name: ObjectServiceNotification.ActionUpdateCurrentUserPortfolioWallet.rawValue,
                                               object: nil)
        NotificationCenter.default.addObserver(self, selector: #selector(updateNetWorth(_:)),
                                               name: ObjectServiceNotification.ActionUpdateCurrentUserPortfolioNetWorth.rawValue,
                                               object: nil)
        NotificationCenter.default.addObserver(self,
                                               selector: #selector(updateStockPrice(_:)),
                                               name: ObjectServiceNotification.ActionUpdateCurrentUserStockPrice.rawValue,
                                               object: nil)
        self.setupData()
    
        
    }
    
    override func viewDidAppear(_ animated: Bool) {
        super.viewDidAppear(true)
        
        self.setupViews()
    }
    
    override func viewDidDisappear(_ animated: Bool) {
        NotificationCenter.default.removeObserver(self, name: ObjectServiceNotification.ActionUpdateCurrentUserPortfolioWallet.rawValue, object: nil)
        NotificationCenter.default.removeObserver(self, name: ObjectServiceNotification.ActionUpdateCurrentUserPortfolioNetWorth.rawValue, object: nil)
        NotificationCenter.default.removeObserver(self, name: ObjectServiceNotification.ActionUpdateCurrentUserStockPrice.rawValue, object: nil)
    }
    
    //MARK: View Setup Functions
    func setupData() {
        
        if ObjectHandler.sharedInstance.netWorth != 0 {
            self.netWorthLabel.text = String(ObjectHandler.sharedInstance.netWorth)
        } else {
            self.netWorthLabel.text = "Loading..."
        }
        
        self.netWorthPercentageChangeLabel.text = ""
        self.netWorthPercentSignLabel.text = ""
        self.netWorthArrowIcon.isHidden = true
        self.netWorthDollarSignLabel.isHidden = true
        
        if ObjectHandler.sharedInstance.netWorth != 0 {
            self.cashAmountLabel.text = String(ObjectHandler.sharedInstance.wallet)
        } else {
            self.cashAmountLabel.text = "Loading..."
        }
        
        if ObjectHandler.sharedInstance.netWorth != 0 && ObjectHandler.sharedInstance.wallet != 0 {
            self.investmentsAmountLabel.text = String(ObjectHandler.sharedInstance.netWorth - ObjectHandler.sharedInstance.wallet)
        } else {
            self.investmentsAmountLabel.text = "Loading..."
        }
    }
    
    func setupViews() {
        UIApplication.shared.statusBarStyle = .lightContent
        
        let gradient = CAGradientLayer.init()
        gradient.colors = [UIColor.init(red: 1.0/20.0, green: 1.0/30.0, blue: 1.0/48.0, alpha: 0.0),
                           UIColor.init(red: 1.0/20.0, green: 1.0/30.0, blue: 1.0/48.0, alpha: 1.0)]
        gradient.startPoint = CGPoint.init(x: blueView.frame.width/2, y: blueView.frame.minY)
        gradient.endPoint = CGPoint.init(x: blueView.frame.width/2, y: blueView.frame.maxY)
        blueView.layer.addSublayer(gradient)
        
        
        
        
        tableView.reloadData()
    }
    
    //MARK: NotificationHandling
    @objc func updateWallet(_ notification:NSNotification) {
        self.cashAmountLabel.text = String(ObjectHandler.sharedInstance.wallet)
        self.investmentsAmountLabel.text = String(ObjectHandler.sharedInstance.netWorth - ObjectHandler.sharedInstance.wallet)
    }
    
    @objc func updateNetWorth(_ notification:NSNotification) {
        self.netWorthLabel.text = String(ObjectHandler.sharedInstance.netWorth)
        self.investmentsAmountLabel.text = String(ObjectHandler.sharedInstance.netWorth - ObjectHandler.sharedInstance.wallet)
    }
    
    @objc func updateLedger(_ notification:NSNotification) {
        //TODO:
    }
    
    @objc func updateStockPrice(_ notification:NSNotification) {
        tableView.reloadData()
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
        return ObjectHandler.sharedInstance.stockArray.count
    }
    
    func tableView(_ tableView: UITableView, cellForRowAt indexPath: IndexPath) -> UITableViewCell {
        let cell = tableView.dequeueReusableCell(withIdentifier: "portfolioTableViewCell", for: indexPath) as? PortfolioTableViewCell
        
        cell?.tickerLabel.text = ObjectHandler.sharedInstance.stockArray[indexPath.row].name  // portfolioArray[indexPath.row].name
        cell?.nameLabel.text = ObjectHandler.sharedInstance.stockArray[indexPath.row].fullname // portfolioArray[indexPath.row].fullname

        if ObjectHandler.sharedInstance.stockArray[indexPath.row].value <= 0 {
            cell?.costLabel.text = ""
        } else {
            cell?.costLabel.text = String(format: "%.2f", ObjectHandler.sharedInstance.stockArray[indexPath.row].value)
        }
        
        if ObjectHandler.sharedInstance.stockArray[indexPath.row].value <= 0 {
            cell?.recordLabel.text = ""
        } else {
            cell?.recordLabel.text = String(format: "%.2f", ObjectHandler.sharedInstance.stockArray[indexPath.row].recordValue)
        }
        
        let amountChanged = ObjectHandler.sharedInstance.stockArray[indexPath.row].amountChanged
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


