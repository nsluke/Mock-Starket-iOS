//
//  MarketViewController.swift
//  Mock Starket
//
//  Created by Luke Solomon on 3/6/18.
//  Copyright © 2018 Luke Solomon. All rights reserved.
//

import UIKit

class MarketViewController: UIViewController {
    
    //Mark: IBOutlets
    @IBOutlet weak var tableView: UITableView!
    
    
    //Mark: Variables
    

    //Mark: View Lifecycle
    override func viewDidLoad() {
        NotificationCenter.default.addObserver(self,
                                               selector: #selector(updateStocks),
                                               name: ObjectServiceNotification.ActionUpdateStockPrice.rawValue,
                                               object: nil)    }
    
    override func viewDidAppear(_ animated: Bool) {
        
    }
    
    override func viewDidDisappear(_ animated: Bool) {
        NotificationCenter.default.removeObserver(self, name: ObjectServiceNotification.ActionUpdateStockPrice.rawValue, object: nil)
    }
    
    @objc func updateStocks() {
        self.tableView.reloadData()
    }
    
    //Mark: Functions
}

extension MarketViewController: UITableViewDelegate, UITableViewDataSource {
    func tableView(_ tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
        return ObjectHandler.sharedInstance.marketArray.count
    }
    
    func tableView(_ tableView: UITableView, cellForRowAt indexPath: IndexPath) -> UITableViewCell {
        let cell = tableView.dequeueReusableCell(withIdentifier: "portfolioTableViewCell", for: indexPath) as? PortfolioTableViewCell
        
        cell?.tickerLabel.text = ObjectHandler.sharedInstance.marketArray[indexPath.row].name
        cell?.nameLabel.text = ObjectHandler.sharedInstance.marketArray[indexPath.row].fullname
        
        if ObjectHandler.sharedInstance.marketArray[indexPath.row].value <= 0 {
            cell?.costLabel.text = ""
        } else {
            cell?.costLabel.text = String(format: "%.2f", ObjectHandler.sharedInstance.marketArray[indexPath.row].value)
        }
        
        if ObjectHandler.sharedInstance.marketArray[indexPath.row].value <= 0 {
            cell?.recordLabel.text = ""
        } else {
            cell?.recordLabel.text = String(format: "%.2f", ObjectHandler.sharedInstance.marketArray[indexPath.row].recordValue)
        }
        
        let amountChanged = ObjectHandler.sharedInstance.marketArray[indexPath.row].amountChanged
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
