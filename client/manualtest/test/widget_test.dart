import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:get/get.dart';
import 'package:manualtest/login.dart';

import 'tap.dart';

void main() {
  testWidgets('Check LoginPage', (WidgetTester tester) async {
    await tester.pumpWidget(
      Material(
        child: GetMaterialApp(
          home: SigIn(),
        ),
      ),
    );

    var login = find.byKey(const Key('Login'));
    expect(login, findsOneWidget);

    await tapWithFindByKey(
      tester,
      const Key('Login'),
    );

    var fehler1 = find.byKey(const Key('Fehler1'));
    expect(fehler1, findsOneWidget);
    var ok1 = find.byKey(const Key('OK1'));
    expect(ok1, findsOneWidget);
    await tapWithFindByKey(
      tester,
      const Key('OK1'),
    );

    var email = find.byKey(const Key('Email'));
    expect(email, findsOneWidget);
    await tester.enterText(email, "email");

    await tapWithFindByKey(
      tester,
      const Key('Login'),
    );
    var fehler2 = find.byKey(const Key('Fehler2'));
    expect(fehler2, findsOneWidget);
    var ok2 = find.byKey(const Key('OK2'));
    expect(ok2, findsOneWidget);
    await tapWithFindByKey(
      tester,
      const Key('OK2'),
    );
  });
}
