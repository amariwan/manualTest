import 'package:flutter/material.dart';
import 'package:fwidget/widgets.dart';
import 'package:get/get.dart';
import 'package:manualtest/stream_controller.dart';

import 'order.dart';

class FTableRowOrdersModel extends FTableRow<OrdersModel> {
  String tester;
  FTableRowOrdersModel({required OrdersModel value, required this.tester})
      : super(value: value);

  static List<String> headers() {
    return [
      'Erzeuger',
      'Beschreibung',
      'Erzeugungsdatum',
      'Ablaufdatum',
      'Status Info',
      'Aktion'
    ];
  }

  @override
  List<Widget> build(BuildContext context) {
    DateTime expireDate = DateTime.parse(value.expireDate);
    String formatExpireDate =
        "${expireDate.year}-${addLeadingZero(expireDate.month)}-${addLeadingZero(expireDate.day)} "
        "${addLeadingZero(expireDate.hour)}:${addLeadingZero(expireDate.minute)}";
    var x = value.createdAt.replaceAll(new RegExp(r" m=\+[\d\.]+"), "");
    DateTime createdAt = DateTime.parse(
        x.split(".")[0].replaceAll(new RegExp(r" m=\+[\d\.]+"), ""));
    String formatCreatedAt =
        "${createdAt.year}-${addLeadingZero(createdAt.month)}-${addLeadingZero(createdAt.day)} "
        "${addLeadingZero(createdAt.hour)}:${addLeadingZero(createdAt.minute)}";
    var bool = expireDate.isBefore(DateTime.now()) || value.testComplete;
    return [
      Text(value.name,
          textAlign: TextAlign.left,
          style: const TextStyle(fontWeight: FontWeight.bold)),
      Text(value.description,
          textAlign: TextAlign.left,
          style: const TextStyle(fontWeight: FontWeight.bold)),
      Text(formatCreatedAt,
          textAlign: TextAlign.left,
          style: const TextStyle(fontWeight: FontWeight.bold)),
      Text(formatExpireDate,
          textAlign: TextAlign.left,
          style: const TextStyle(fontWeight: FontWeight.bold)),
      Text(
          textAlign: TextAlign.left,
          bool ? 'Abgeschlossen' : 'Aktive',
          style: const TextStyle(fontWeight: FontWeight.bold)),
      Wrap(
        spacing: 8.0,
        alignment: WrapAlignment.start,
        children: [
          TextButton(
              onPressed: () => Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (context) => OrderPage(
                        tester: tester,
                        title: 'Orders',
                        orderId: value.id,
                      ),
                    ),
                  ),
              child: const Text('Öffnen')),
          if (DateTime.parse(value.expireDate).isBefore(DateTime.now()) ||
              value.testComplete)
            TextButton(
                onPressed: () {
                  Get.find<WebSocketController>().deleteValues(value.id);
                },
                child: const Text('Löschen')),
        ],
      ),
    ];
  }

  @override
  void updateValue(OrdersModel newValue) {
    value = newValue;
  }

  @override
  bool hasSameValue(FTableRow<OrdersModel> other) {
    return value == other.value;
  }

  @override
  FTableRow<OrdersModel> clone() {
    return FTableRowOrdersModel(
      value: value.clone(),
      tester: tester,
    );
  }

  //@override
  //int compareCell(int i, FTableRow<ConvCustomConverterModel> other) {
  //  switch (i) {
  //    case 0:
  //      return Comparable.compare(value.from?.value ?? '',
  //          (other as FTableRowModel).value.from?.value ?? '');
  //    case 1:
  //      return Comparable.compare(value.to?.value ?? '',
  //          (other as FTableRowModel).value.to?.value ?? '');
  //    case 2:
  //      return Comparable.compare(
  //          value.code.value, (other as FTableRowModel).value.code.value);
  //  }
  //  return 0;
  //}
}

class OrdersModel implements FCloneable<OrdersModel> {
  String name;
  String description;
  String createdAt;
  String expireDate;
  String id;
  bool testComplete;
  OrdersModel(
      {required this.name,
      required this.createdAt,
      required this.description,
      required this.expireDate,
      required this.id,
      required this.testComplete});

  @override
  OrdersModel clone() {
    return OrdersModel(
      name: name,
      expireDate: expireDate,
      createdAt: createdAt,
      description: description,
      id: id,
      testComplete: testComplete,
    );
  }

  @override
  bool operator ==(Object other) {
    if (other is OrdersModel) {
      return id == other.id &&
          testComplete == other.testComplete &&
          name == other.name &&
          expireDate == other.expireDate &&
          createdAt == other.createdAt &&
          description == other.description;
    }
    return false;
  }

  @override
  int get hashCode {
    return Object.hash(
      name,
      expireDate,
      createdAt,
      description,
      id,
      testComplete,
    );
  }
}

String addLeadingZero(int value) {
  if (value < 10) {
    return "0$value";
  } else {
    return "$value";
  }
}
