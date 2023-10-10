import 'package:flutter/material.dart';
import 'package:fwidget/widgets.dart';
import 'package:get/get.dart';
import 'package:manualtest/stream_controller.dart';

import 'json.dart';

class FTableRowOrderModel extends FTableRow<OrderModel> {
  void Function()? onChange;
  final Tasks task;
  String id;
  String tester;
  FTableRowOrderModel(
      {required OrderModel value,
      this.onChange,
      required this.id,
      required this.tester,
      required this.task})
      : super(value: value) {
    isTacked = task.tester == tester;
  }
  bool isTacked = false;
  static List<String> headers() {
    return ['Tester', 'Beschreibung', 'Status', 'Kommentar', 'Aktion'];
  }

  @override
  List<Widget> build(BuildContext context) {
    value.commitController.selection = TextSelection(
        baseOffset: value.commitController.text.length,
        extentOffset: value.commitController.text.length);
    return [
      value.task.tester == ''
          ? TextButton(
              onPressed: () {
                if (onChange != null) {
                  value.task.tester = tester;
                  isTacked = true;
                  onChange!();
                }
                //setState(() {
                //  widget.orderModels[id].tester = widget.tester;
                //});
              },
              child: const Text("Ich übernehme",
                  textAlign: TextAlign.center,
                  style: TextStyle(fontWeight: FontWeight.bold)))
          : FInputText(
              labelText: '',
              controller: FTextEditingController(text: value.task.tester),
              readOnly: true,
            ),
      Text(value.description,
          textAlign: TextAlign.left,
          style: const TextStyle(fontWeight: FontWeight.bold)),
      IgnorePointer(
        ignoring: value.task.tester != tester && value.task.tester != "",
        child: Wrap(
          spacing: 8.0,
          alignment: WrapAlignment.start,
          children: [
            value.task.result == null || value.task.result == ''
                ? TextButton(
                    onPressed: () {
                      if (onChange != null) {
                        value.task.result = 'success';
                        if (value.task.tester != tester) {
                          value.task.tester = tester;
                        }
                        onChange!();
                      }
                    },
                    child: const Text(
                      "Erfolgreich",
                      textAlign: TextAlign.center,
                      style: TextStyle(fontWeight: FontWeight.bold),
                    ),
                  )
                : value.task.result == 'success'
                    ? const Icon(
                        Icons.check_circle_outline,
                        size: 40,
                        color: Colors.green,
                      )
                    : const Text(""),
            value.task.result == null || value.task.result == ''
                ? TextButton(
                    onPressed: () {
                      if (onChange != null) {
                        value.task.result = 'failed';
                        if (value.task.tester != tester) {
                          value.task.tester = tester;
                        }
                        onChange!();
                      }
                    },
                    child: const Text(
                      "Fehlerhaft",
                      textAlign: TextAlign.center,
                      style: TextStyle(fontWeight: FontWeight.bold),
                    ),
                  )
                : value.task.result == 'failed'
                    ? const Icon(
                        Icons.clear,
                        size:
                            40, // Hier können Sie die Größe des Icons anpassen
                        color: Colors
                            .red, // Hier können Sie die Farbe des Icons anpassen
                      )
                    : const Text(""),
          ],
        ),
      ),
      IgnorePointer(
        ignoring: value.task.tester != tester && value.task.tester != "",
        child: Container(
          alignment: Alignment.topLeft,
          child: FInputText(
            controller: value.commitController,
            onChange: () {
              if (onChange != null) {
                value.task.comment = value.commitController.text;
                if (value.task.tester != tester) {
                  value.task.tester = tester;
                }
                onChange!();
              }
            },
            labelText: null,
          ),
        ),
      ),
      Container(
        alignment: Alignment.topLeft,
        child: TextButton(
            onPressed: () {
              if (onChange != null) {
                isTacked = false;
                if (value.task.tester != tester) {
                  Get.find<WebSocketController>()
                      .sendEmail(id, value.task.tester, value.task.description);
                }
                value.task.tester = '';
                value.task.result = '';
                value.task.comment = '';
                onChange!();
              }
            },
            child: const Text("zurücksetzen",
                textAlign: TextAlign.center,
                style: TextStyle(fontWeight: FontWeight.bold))),
      ),
    ];
  }

  @override
  void updateValue(OrderModel newValue) {
    value = newValue;
  }

  @override
  bool hasSameValue(FTableRow<OrderModel> other) {
    return value == other.value;
  }

  @override
  FTableRow<OrderModel> clone() {
    return FTableRowOrderModel(
      id: id,
      task: task,
      onChange: onChange,
      tester: tester,
      value: value.clone(),
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

class OrderModel extends FCloneable<OrderModel> {
  final FocusNode focusNode = FocusNode(debugLabel: 'OrderedListCheckboxModel');
  final Tasks task;
  final FTextEditingController testerController;
  final FTextEditingController descriptionController;
  bool? status;
  final FTextEditingController commitController;

  OrderModel({required this.task})
      : descriptionController = FTextEditingController(text: task.description),
        testerController = FTextEditingController(text: task.tester),
        commitController = FTextEditingController(text: task.comment);

  String get tester {
    return testerController.text;
  }

  String get description {
    return descriptionController.text;
  }

  String get commit {
    return commitController.text;
  }

  set tester(String value) {
    testerController.text = value;
  }

  set description(String value) {
    descriptionController.text = value;
  }

  set commit(String value) {
    commitController.text = value;
  }

  bool? get statusValue {
    return status;
  }

  set statusValue(bool? value) {
    status = value;
  }

  bool isEqual(OrderModel other) {
    if (tester != other.tester ||
        commit != other.commit ||
        description != other.description ||
        statusValue != other.statusValue) {
      return false;
    }
    return true;
  }

  @override
  bool operator ==(Object other) {
    if (other is OrderModel) {
      return isEqual(other);
    }
    return false;
  }

  @override
  OrderModel clone() {
    return OrderModel(task: task);
  }

  @override
  int get hashCode => Object.hash(
        statusValue,
        commit,
        description,
        tester,
      );
}
