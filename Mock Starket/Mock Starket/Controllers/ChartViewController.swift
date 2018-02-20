//
//  ChartViewController.swift
//  Mock Starket
//
//  Created by Luke Solomon on 1/21/18.
//  Copyright © 2018 Luke Solomon. All rights reserved.
//

import UIKit
import Starscream
import Charts

class ChartViewController: UIViewController {
    
    var socket = WebSocket(url: URL(string: "ws://159.89.154.221:8000/ws")!)
    
    @IBOutlet weak var beginButton: UIButton!
    @IBOutlet weak var lineChart: LineChartView!
    
    
    // MARK: View lifeCycle
    override func viewDidLoad() {
        super.viewDidLoad()
        
        // establishing connection to the websocket
        socket.delegate = self
        
        let dataPoints = ["Daenerys", "Jon Snow", "Rhaegar Targaryen"]
        let values = [22.0, 25.0, 29.0]
        setChart(dataPoints: dataPoints, values: values)
    }
    
    override func viewWillAppear(_ animated: Bool) {
        

    }
    
    override func viewDidAppear(_ animated: Bool) {
        
    }
    
    func setChart(dataPoints: [String], values: [Double]) {
        
        
        var dataEntries: [ChartDataEntry] = []
        
        for i in 0..<dataPoints.count {
            let dataEntry = ChartDataEntry(x: values[i], y: Double(i))
            dataEntries.append(dataEntry)
        }
        
        let lineChartDataSet =  LineChartDataSet(values: dataEntries, label: "prices")  // LineChartDataSet(values: dataEntries, label: "Units Sold")
        let lineChartData = LineChartData(dataSet: lineChartDataSet)
        
//        LineChartData(values: dataPoints, dataSet: lineChartDataSet)
//        LineChartData(dataSets: [IChartDataSet]?)
        
        lineChart.data = lineChartData
        
    }
    
    deinit {
        socket.disconnect(forceTimeout: 0)
        socket.delegate = nil
    }
    
    func sendMessage(_ message: String) {
        socket.write(string: message)
        
    }
    
    @IBAction func connect(_ sender: Any) {
        socket.connect()
        
    }
}

extension ChartViewController: WebSocketDelegate {
    
    func websocketDidConnect(socket: WebSocketClient) {
        print("Connected")

        socket.write(string: "{\"action\": \"login\", \"value\": {\"username\": \"username\", \"password\":\"password\"}}")
    }
    
    func websocketDidDisconnect(socket: WebSocketClient, error: Error?) {
        
    }
    
    func websocketDidReceiveMessage(socket: WebSocketClient, text: String) {
        print(text)
    }
    
    func websocketDidReceiveData(socket: WebSocketClient, data: Data) {
        print(data)
    }
}
