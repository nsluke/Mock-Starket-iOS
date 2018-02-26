//
//  PortfolioTableViewCell.swift
//  Mock Starket
//
//  Created by Luke Solomon on 2/20/18.
//  Copyright © 2018 Luke Solomon. All rights reserved.
//

import UIKit

class PortfolioTableViewCell: UITableViewCell {

    @IBOutlet weak var tickerLabel: UILabel!
    @IBOutlet weak var nameLabel: UILabel!
    @IBOutlet weak var recordLabel: UILabel!
    @IBOutlet weak var changeLabel: UILabel!
    @IBOutlet weak var costLabel: UILabel!
    
    
    override func awakeFromNib() {
        super.awakeFromNib()
        // Initialization code
    }

    override func setSelected(_ selected: Bool, animated: Bool) {
        super.setSelected(selected, animated: animated)

        // Configure the view for the selected state
    }

}
