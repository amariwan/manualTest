import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:fwidget/widgets.dart';

import 'json.dart';
import 'order_model.dart';
import 'orders.dart';
import 'stream_controller.dart';

class OrderPage extends StatefulWidget {
  final String title;
  final String tester;
  final String orderId;
  OrderPage({
    Key? key,
    required this.title,
    required this.tester,
    required this.orderId,
  }) : super(key: key) {
    Get.find<WebSocketController>().getValues();
  }

  @override
  _OrderPageState createState() => _OrderPageState();
  bool buildFirstTime = true;
  List<OrderModel> orderModels = [];
}

class _OrderPageState extends State<OrderPage> {
  _OrderPageState();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        automaticallyImplyLeading: false,
        leading: Row(
          children: [
            IconButton(
                icon: const Icon(Icons.arrow_back),
                onPressed: () => Navigator.push(
                      context,
                      MaterialPageRoute(
                        builder: (context) => OrdersPages(
                          tester: widget.tester,
                          title: 'Orders',
                        ),
                      ),
                    )),
          ],
        ),
        title: Row(
          children: [
            Expanded(
                flex: 1,
                child: Container(
                    alignment: Alignment.topLeft,
                    child: Image.asset('assets/images/logo.png'))),
            const Expanded(flex: 1, child: Text('Test')),
          ],
        ),
        backgroundColor: Colors.red,
      ),
      body: Container(
        constraints: const BoxConstraints.expand(),
        decoration: const BoxDecoration(
            image: DecorationImage(image: NetworkImage(""), fit: BoxFit.cover)),
        child: Padding(
          padding: const EdgeInsets.all(20.0),
          child: StreamBuilder(
            stream: Get.find<WebSocketController>().stream,
            builder: (context, snapshot) {
              List<Order> orders = _convertInfoFromWebSocketOrder(
                snapshot,
                widget.orderId,
              );
              return Padding(
                padding: const EdgeInsets.symmetric(vertical: 24.0),
                child: ListView(
                  children: <Widget>[
                    if (snapshot.hasData)
                      getWidgetsFromWebsocketOrder(
                          orders, widget.orderId, widget.orderModels),
                  ],
                ),
              );
            },
          ),
        ),
      ),
    );
  }

  Widget getWidgetsFromWebsocketOrder(List<Order> infoFromWebSocket,
      String orderId, List<OrderModel> orderModels) {
    bool isActive = true;
    List<List<OrderModel>> models = [];
    for (int i = 0; i < infoFromWebSocket.length; i++) {
      List<OrderModel> iModel = [];
      for (int j = 0; j < infoFromWebSocket[i].tasks.length; j++) {
        iModel.add(OrderModel(
          task: infoFromWebSocket[i].tasks[j],
        ));
      }
      models.add(iModel);
    }
    FTableController<OrderModel> controller;
    List<FTableRowOrderModel> tableRows = [];
    for (int i = 0; i < infoFromWebSocket.length; i++) {
      if (infoFromWebSocket[i].id == orderId) {
        isActive = !infoFromWebSocket[i].testComplete;
        for (int j = 0; j < infoFromWebSocket[i].tasks.length; j++) {
          tableRows.add(
            FTableRowOrderModel(
                id: infoFromWebSocket[i].id,
                task: infoFromWebSocket[i].tasks[j],
                tester: widget.tester,
                onChange: () {
                  List<Tasks> tasks = [];
                  Order infoToSendToServer;
                  Map<String, dynamic> jsonToSend = {};
                  for (int l = 0; l < infoFromWebSocket.length; l++) {
                    if (infoFromWebSocket[l].id == orderId) {
                      for (int k = 0;
                          k < infoFromWebSocket[l].tasks.length;
                          k++) {
                        tasks.add(
                          Tasks(
                              result: models[l][k].task.result,
                              comment: models[l][k].task.comment,
                              description: models[l][k].task.description,
                              tester: models[l][k].task.tester),
                        );
                        DateTime date =
                            DateTime.parse(infoFromWebSocket[l].expireDate);
                        if (date.isBefore(DateTime.now()) ||
                            infoFromWebSocket[l].testComplete) {
                          isActive = false;
                        }
                        infoToSendToServer = Order(
                            event: "updateValue",
                            callbackUrl: infoFromWebSocket[l].callbackUrl,
                            creator: infoFromWebSocket[l].creator,
                            description: infoFromWebSocket[l].description,
                            expireDate: infoFromWebSocket[l].expireDate,
                            mail: infoFromWebSocket[l].mail,
                            tasks: tasks,
                            testers: infoFromWebSocket[l].testers,
                            testComplete: infoFromWebSocket[l].testComplete,
                            id: infoFromWebSocket[l].id,
                            createdAt: infoFromWebSocket[l].createdAt,
                            updatedAt: infoFromWebSocket[l].updatedAt,
                            stump: infoFromWebSocket[l].stump);
                        jsonToSend = infoToSendToServer.toJson();
                      }
                    }
                  }
                  setState(() {
                    Get.find<WebSocketController>()
                        .sendMessage(json.encode(jsonToSend));
                  });
                },
                value: models[i][j]),
          );
        }
      }
    }
    Order toJenkins = Order(
        event: "getValues",
        testComplete: true,
        callbackUrl: '',
        createdAt: '',
        creator: Tester(email: '', name: ''),
        description: '',
        expireDate: '',
        id: '',
        mail: Mail(body: '', subject: ''),
        tasks: [],
        testers: [],
        updatedAt: TimeOfDay.now().toString(),
        stump: Signature(
            index: Get.find<Index>().index,
            createdAt: TimeOfDay.now().toString()));
    Get.find<Index>().increment();

    controller = FTableController(rows: tableRows);

    bool completed = true;
    for (int l = 0; l < infoFromWebSocket.length; l++) {
      if (infoFromWebSocket[l].id == orderId) {
        for (int i = 0; i < infoFromWebSocket[l].tasks.length; i++) {
          if (infoFromWebSocket[l].tasks[i].result == "" ||
              infoFromWebSocket[l].tasks[i].result == null) {
            completed = false;
          }
        }
      }
    }

    var output = Center(
      child: IgnorePointer(
        ignoring: !isActive,
        child: Column(
          children: [
            FTable(
              controller: controller,
              headers: FTableRowOrderModel.headers(),
              columnWidths: const [
                IntrinsicColumnWidth(),
                FixedColumnWidth(300),
                FixedColumnWidth(200),
                FlexColumnWidth(150),
                FlexColumnWidth(50),
              ],
              //onChange: widget.onChange,
            ),
            for (int l = 0; l < infoFromWebSocket.length; l++)
              if (infoFromWebSocket[l].id == orderId &&
                  completed &&
                  !infoFromWebSocket[l].testComplete)
                TextButton(
                  onPressed: () {
                    infoFromWebSocket[l].testComplete = true;

                    toJenkins = Order(
                        event: "protokol",
                        testComplete: infoFromWebSocket[l].testComplete,
                        callbackUrl: infoFromWebSocket[l].callbackUrl,
                        createdAt: infoFromWebSocket[l].createdAt,
                        creator: infoFromWebSocket[l].creator,
                        description: infoFromWebSocket[l].description,
                        expireDate: infoFromWebSocket[l].expireDate,
                        id: infoFromWebSocket[l].id,
                        mail: infoFromWebSocket[l].mail,
                        tasks: infoFromWebSocket[l].tasks,
                        testers: infoFromWebSocket[l].testers,
                        updatedAt: TimeOfDay.now().toString(),
                        stump: Signature(
                            index: Get.find<Index>().index,
                            createdAt: TimeOfDay.now().toString()));
                    Get.find<Index>().increment();

                    setState(() {
                      Get.find<WebSocketController>()
                          .sendMessage(json.encode(toJenkins.toJson()));
                    });
                  },
                  child: const Text(
                    "Protokoll senden",
                    textAlign: TextAlign.center,
                    style: TextStyle(fontWeight: FontWeight.bold),
                  ),
                )
          ],
        ),
      ),
    );
    return output;
  }

  List<Order> _convertInfoFromWebSocketOrder(
      AsyncSnapshot<dynamic> snapshot, String orderId) {
    List<Order> infoFromWebSocket = [];
    if (snapshot.hasData) {
      var overview = Orders.fromJson(json.decode(snapshot.data));
      print(overview.stump.index.toString() +
          overview.stump.createdAt.toString());
      for (var element in overview.orders) {
        infoFromWebSocket.add(
          Order(
              event: element.event,
              callbackUrl: element.callbackUrl,
              creator: element.creator,
              description: element.description,
              expireDate: element.expireDate,
              mail: element.mail,
              tasks: element.tasks,
              testers: element.testers,
              testComplete: element.testComplete,
              id: element.id,
              createdAt: element.createdAt,
              updatedAt: TimeOfDay.now().toString(),
              stump: Signature(
                  index: Get.find<Index>().index,
                  createdAt: TimeOfDay.now().toString())),
        );
        Get.find<Index>().increment();
        if (element.id == orderId) {
          if (widget.buildFirstTime) {
            for (var task in element.tasks) {
              widget.orderModels.add(OrderModel(task: task));
            }
            widget.buildFirstTime = false;
          }
        }
      }
    }

    return infoFromWebSocket;
  }
}
