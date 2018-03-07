//
//  MarketViewController.swift
//  Mock Starket
//
//  Created by Luke Solomon on 3/6/18.
//  Copyright © 2018 Luke Solomon. All rights reserved.
//

import UIKit



class MarketViewController: ViewController {
    
    //Mark: IBOutlets
    @IBOutlet weak var tableView: UITableView!
    
    
    //Mark: Variables
    var stockArray = [
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
    var mutableStockSet = NSMutableOrderedSet()

    //Mark: View Lifecycle
    override func viewDidLoad() {
        NotificationCenter.default.addObserver(self, selector: #selector(update(_:)), name: NetworkServiceNotification.SocketMessageReceived.rawValue, object: nil)
    }
    
    override func viewDidAppear(_ animated: Bool) {
        for stock in stockArray {
            mutableStockSet.add(stock.name)
        }
    }
    
    override func viewDidDisappear(_ animated: Bool) {
        NotificationCenter.default.removeObserver(self, name: NetworkServiceNotification.SocketMessageReceived.rawValue, object: nil)
    }
    
    //Mark: Functions
    
    
    @objc func update(_ notification:NSNotification) {
        guard let actionArray = notification.userInfo?["actionArray"] as? [Action] else {
            return
        }
        
        for action in actionArray {
            if action.type == "update" {
                if let valuable = action.value.valuable {
                    let stock = Stock.init(name: valuable.tickerID, value: valuable.current_price)
                    let index = mutableStockSet.index(of: valuable.tickerID)
                    var amountChanged = 0.0
                    
                    if valuable.current_price != 0 && stockArray[index].value != 0 {
                        amountChanged = valuable.current_price - stockArray[index].value
                    }
                    
                    if mutableStockSet.contains(stock.name) {
                        stockArray.remove(at: index)
                        stockArray.insert(stock, at: index)
                        if stockArray[index].recordValue < stock.value {
                            stockArray[index].recordValue = stock.value
                        }
                        
                        stockArray[index].amountChanged = amountChanged
                    } else {
                        mutableStockSet.add(stock.name)
                        stockArray.append(stock)
                        print("New Field!" + stock.name)
                    }
                }

                if let portfolio = action.value.portfolio {
                    
                }
                
                self.tableView.reloadData()

            }
        }
    }
    
}
extension MarketViewController: UITableViewDelegate, UITableViewDataSource {
    func tableView(_ tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
        return stockArray.count
    }
    
    func tableView(_ tableView: UITableView, cellForRowAt indexPath: IndexPath) -> UITableViewCell {
        let cell = tableView.dequeueReusableCell(withIdentifier: "portfolioTableViewCell", for: indexPath) as? PortfolioTableViewCell
        
        cell?.tickerLabel.text = stockArray[indexPath.row].name
        cell?.nameLabel.text = stockArray[indexPath.row].fullname
        
        if stockArray[indexPath.row].value <= 0 {
            cell?.costLabel.text = ""
        } else {
            cell?.costLabel.text = String(format: "%.2f", stockArray[indexPath.row].value)
        }
        
        if stockArray[indexPath.row].value <= 0 {
            cell?.recordLabel.text = ""
        } else {
            cell?.recordLabel.text = String(format: "%.2f", stockArray[indexPath.row].recordValue)
        }
        
        let amountChanged = stockArray[indexPath.row].amountChanged
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
