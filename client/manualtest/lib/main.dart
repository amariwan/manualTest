import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:manualtest/stream_controller.dart';

import 'login.dart';
import 'order.dart';

/// `runApp` takes a `Widget` and displays it in the screen
void main() => runApp(const ManualTest());

/// It creates a stateful widget.
class ManualTest extends StatefulWidget {
  const ManualTest({super.key});

  @override
  State<ManualTest> createState() => _ManualTestState();
}

/// > This class is a stateful widget that has a state class called _ManualTestState
class _ManualTestState extends State<ManualTest> {
  @override
  void initState() {
    super.initState();

    /// Creating a new instance of the WebSocketController class.
    WebSocketController cannel = WebSocketController('ws://localhost:3000/ws');
    /// `Get.put(cannel);` is registering an instance of the `WebSocketController` class with the GetX
    /// dependency injection system. This allows the instance to be easily accessed and used throughout
    /// the app by calling `Get.find<WebSocketController>()`.
    Get.put(cannel);

    Get.put(Index(0));
  }

  @override
  void dispose() {
    super.dispose();
    Get.find<WebSocketController>().dispose();
  }

  @override

  /// A function that returns a widget.
  ///
  /// Args:
  ///   context (BuildContext): The context is used to access the theme and localizations of the app.
  Widget build(BuildContext context) {
    var id = Uri.base.queryParameters["id"] ?? "";
    var tester = Uri.base.queryParameters["tester"] ?? "";
    if (Uri.base.toString().contains('id')) {
      return GetMaterialApp(
        title: 'title',
        home: OrderPage(
          tester: tester,
          title: 'Orders',
          orderId: id,
        ),
      );
    }
    const title = 'ManualTest';
    return GetMaterialApp(
      title: title,
      home: SigIn(),
    );
  }
}
