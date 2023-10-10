import 'package:flutter/material.dart';

class Tester {
  String email = '';
  String name = '';

  Tester({required this.email, required this.name});

  Tester.fromJson(Map<String, dynamic> json) {
    email = json['email'] ?? '';
    name = json['name'] ?? '';
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['email'] = email;
    data['name'] = name;
    return data;
  }
}

class Signature {
  int index = 0;
  String createdAt = '';

  Signature({required this.index, required this.createdAt});

  Signature.fromJson(Map<String, dynamic> json) {
    index = json['index'];
    createdAt = json['createdAt'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['index'] = index;
    data['createdAt'] = createdAt;
    return data;
  }
}

class Mail {
  String body = '';
  String subject = '';

  Mail({required this.body, required this.subject});

  Mail.fromJson(Map<String, dynamic> json) {
    body = json['body'];
    subject = json['subject'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['body'] = body;
    data['subject'] = subject;
    return data;
  }
}

class Tasks {
  String description = '';
  String tester = '';
  String? result;
  String? comment = '';

  Tasks(
      {required this.description,
      required this.tester,
      this.result,
      this.comment});

  Tasks.fromJson(Map<String, dynamic> json) {
    description = json['description'] ?? '';
    tester = json['tester'] ?? '';
    result = json['result'];
    comment = json['comment'] ?? '';
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['description'] = description;
    data['tester'] = tester;
    data['result'] = result;
    data['comment'] = comment;
    return data;
  }
}

class Order {
  String event = '';
  String callbackUrl = '';
  Tester creator = Tester(email: '', name: '');
  String description = '';
  String expireDate = '';
  Mail mail = Mail(body: '', subject: '');
  List<Tasks> tasks = [];
  List<Tester> testers = [];
  bool testComplete = false;
  String id = '';
  String createdAt = TimeOfDay.now().toString();
  String updatedAt = TimeOfDay.now().toString();
  Signature stump = Signature(index: 0, createdAt: TimeOfDay.now().toString());

  Order(
      {required this.event,
      required this.callbackUrl,
      required this.creator,
      required this.description,
      required this.expireDate,
      required this.mail,
      required this.tasks,
      required this.testers,
      required this.testComplete,
      required this.id,
      required this.createdAt,
      required this.updatedAt,
      required this.stump});

  Order.fromJson(Map<String, dynamic> json) {
    event = json['event'];
    callbackUrl = json['callbackUrl'];
    creator = json['creator'] != null
        ? Tester.fromJson(json['creator'])
        : Tester(email: '', name: '');
    description = json['description'];
    expireDate = json['expireDate'];
    mail = json['mail'] != null
        ? Mail.fromJson(json['mail'])
        : Mail(body: '', subject: '');
    if (json['tasks'] != null) {
      tasks = [];
      json['tasks'].forEach((v) {
        tasks.add(Tasks.fromJson(v));
      });
    }
    if (json['testers'] != null) {
      testers = [];
      json['testers'].forEach((v) {
        testers.add(Tester.fromJson(v));
      });
    }
    id = json['id'];
    createdAt = json['createdAt'];
    testComplete = json['testComplete'];
    updatedAt = json['updatedAt'];

    stump = json['stump'] != null
        ? Signature.fromJson(json['stump'])
        : Signature(index: 0, createdAt: TimeOfDay.now().toString());
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['event'] = event;
    data['callbackUrl'] = callbackUrl;

    data['creator'] = creator.toJson();

    data['description'] = description;
    data['expireDate'] = expireDate;

    data['mail'] = mail.toJson();

    data['tasks'] = tasks.map((v) => v.toJson()).toList();

    data['testers'] = testers.map((v) => v.toJson()).toList();

    data['id'] = id;
    data['createdAt'] = createdAt;
    data['testComplete'] = testComplete;
    data['updatedAt'] = updatedAt;
    data['stump'] = stump.toJson();
    return data;
  }
}

class Orders {
  List<Order> orders = [];
  String event = '';
  Signature stump = Signature(index: 0, createdAt: TimeOfDay.now().toString());
  Orders({required this.orders, required this.event, required this.stump});

  Orders.fromJson(Map<String, dynamic> json) {
    event = json['event'];
    if (json['orders'] != null) {
      orders = [];
      json['orders'].forEach((v) {
        orders.add(Order.fromJson(v));
      });
    }
    stump = json['stump'] != null
        ? Signature.fromJson(json['stump'])
        : Signature(index: 0, createdAt: TimeOfDay.now().toString());
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['event'] = event;
    data['orders'] = orders;
    data['stump'] = stump;
    return data;
  }
}
