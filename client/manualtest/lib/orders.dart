import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:fwidget/widgets.dart';
import 'package:get/get.dart';
import 'package:manualtest/stream_controller.dart';

import 'json.dart';
import 'login.dart';
import 'orders_model.dart';

class OrdersPages extends StatefulWidget {
  final String title;
  final String tester;

  OrdersPages({Key? key, required this.title, required this.tester})
      : super(key: key) {
    Get.find<WebSocketController>().getValues();
  }

  @override
  _OrdersPagesState createState() => _OrdersPagesState();
}

class _OrdersPagesState extends State<OrdersPages> {
  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        automaticallyImplyLeading: false,
        leading: IconButton(
            icon: const Icon(Icons.arrow_back),
            onPressed: () => Navigator.push(
                  context,
                  MaterialPageRoute(
                    builder: (context) => SigIn(),
                  ),
                )),
        title: Row(
          children: [
            Expanded(
                flex: 1,
                child: Container(
                    alignment: Alignment.topLeft,
                    child: Image.asset('assets/images/logo.png'))),
            const Expanded(flex: 1, child: Text('Test Hosts')),
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
            // ignore: missing_return
            builder: (context, snapshot) {
              List<Order> infoFromWebSocket =
                  _convertInfoFromWebSocket(snapshot);
              return Padding(
                padding: const EdgeInsets.symmetric(vertical: 24.0),
                child: ListView(
                  children: <Widget>[
                    if (snapshot.hasData)
                      getWidgetsFromWebsocket(infoFromWebSocket),
                  ],
                ),
              );
            },
          ),
        ),
      ),
    );
  }

  Widget getWidgetsFromWebsocket(List<Order> infoFromWebSocket) {
    FTableController<OrdersModel> controller;

    List<FTableRowOrdersModel> tableRows = [];
    for (var element in infoFromWebSocket) {
      tableRows.add(
        FTableRowOrdersModel(
          value: OrdersModel(
            name: element.creator.name,
            createdAt: element.createdAt,
            description: element.description,
            expireDate: element.expireDate,
            id: element.id,
            testComplete: element.testComplete,
          ),
          tester: widget.tester,
        ),
      );
    }

    controller = FTableController(rows: tableRows);
    var output = Center(
      child: FTable(
        controller: controller,
        headers: FTableRowOrdersModel.headers(),
        columnWidths: const [
          IntrinsicColumnWidth(),
          FlexColumnWidth(200),
          FlexColumnWidth(100),
          FlexColumnWidth(100),
          FlexColumnWidth(50),
          FlexColumnWidth(100),
        ],
        //onChange: widget.onChange,
      ),
    );

    return output;
  }

  List<Order> _convertInfoFromWebSocket(AsyncSnapshot<dynamic> snapshot) {
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
              updatedAt: element.updatedAt,
              stump: element.stump),
        );
      }
    }
    return infoFromWebSocket;
  }
}
