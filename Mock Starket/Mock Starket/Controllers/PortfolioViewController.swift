//
//  PortfolioViewController.swift
//  Mock Starket
//
//  Created by Luke Solomon on 2/20/18.
//  Copyright © 2018 Luke Solomon. All rights reserved.
//

import UIKit
import SideMenu

class PortfolioViewController: UIViewController {
    
    @IBOutlet weak var headerView: UIView!
    @IBOutlet weak var sideMenuButton: UIButton!
    @IBOutlet weak var blueView: UIView!
    @IBOutlet weak var tableViewHeaderView: UIView!
    @IBOutlet weak var portfolioStampView: UIView!
    
    
    override func viewDidLoad() {
        super.viewDidLoad()
    }
    
    override func viewDidAppear(_ animated: Bool) {
        super.viewDidAppear(true)
        self.setupHeader()
    }
    
    func setupHeader() {
        let gradient = CAGradientLayer.init()
        gradient.colors = [UIColor.init(red: 1.0/20.0, green: 1.0/30.0, blue: 1.0/48.0, alpha: 0.0), 
                           UIColor.init(red: 1.0/20.0, green: 1.0/30.0, blue: 1.0/48.0, alpha: 1.0) ]
        gradient.startPoint = CGPoint.init(x: blueView.frame.width/2, y: blueView.frame.minY)
        gradient.endPoint = CGPoint.init(x: blueView.frame.width/2, y: blueView.frame.minY)
        blueView.layer.addSublayer(gradient)
        
    }
    
    
    //handle button tap
    @IBAction func sideMenuButtonTapped(_ sender: Any) {
        present(SideMenuManager.default.menuLeftNavigationController!, animated: true, completion: nil)
        dismiss(animated: true, completion: nil)
    }
    
}

extension PortfolioViewController: UITabBarDelegate {
    
    func tabBar(_ tabBar: UITabBar, didSelect item: UITabBarItem) {
        
    }

}
