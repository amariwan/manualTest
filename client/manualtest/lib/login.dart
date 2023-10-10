import 'package:flutter/material.dart';

import 'orders.dart';

class SigIn extends StatefulWidget {
  SigIn({
    Key? key,
  }) : super(key: key);

  @override
  _SigInState createState() => _SigInState();
  TextEditingController loginController = TextEditingController();
}

class _SigInState extends State<SigIn> {
  _SigInState();

  @override
  Widget build(BuildContcontext) {
    return Scaffold(
      appBar: AppBar(
        automaticallyImplyLeading: false,
        title: Row(
          children: [
            Expanded(
                flex: 1,
                child: Container(
                    alignment: Alignment.topLeft,
                    child: Image.asset('assets/images/logo.png'))),
            const Expanded(flex: 1, child: Text('Login')),
          ],
        ),
        backgroundColor: Colors.red,
      ),
      body: Padding(
        padding: const EdgeInsets.all(8.0),
        child: Center(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              const SizedBox(height: 100),
              Padding(
                padding: const EdgeInsets.all(8.0),
                child: SizedBox(
                  width: 500,
                  child: TextField(
                    key: const Key('Email'),
                    decoration: const InputDecoration(
                      labelText: 'Email der Tester',
                      border: OutlineInputBorder(),
                    ),
                    controller: widget.loginController,
                  ),
                ),
              ),
              const SizedBox(height: 20),
              TextButton(
                child: const Text('Login', key: Key('Login')),
                onPressed: () {
                  String email = widget.loginController.text;
                  if (email.isEmpty) {
                    showDialog(
                      context: context,
                      builder: (BuildContext context) {
                        return AlertDialog(
                          title: const Text('Fehler'),
                          content: const Text(
                              'Bitte geben Sie Ihre E-Mail-Adresse ein.',
                              key: Key('Fehler1')),
                          actions: [
                            TextButton(
                              child: const Text('OK', key: Key("OK1")),
                              onPressed: () {
                                Navigator.of(context).pop();
                              },
                            ),
                          ],
                        );
                      },
                    );
                  } else if (!RegExp(
                          r'\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b')
                      .hasMatch(email)) {
                    showDialog(
                      context: context,
                      builder: (BuildContext context) {
                        return AlertDialog(
                          title: const Text('Fehler'),
                          content: const Text(
                              'Bitte geben Sie eine gültige E-Mail-Adresse ein.',
                              key: Key('Fehler2')),
                          actions: [
                            TextButton(
                              child: const Text('OK', key: Key('OK2')),
                              onPressed: () {
                                Navigator.of(context).pop();
                              },
                            ),
                          ],
                        );
                      },
                    );
                  } else {
                    Navigator.push(
                      context,
                      MaterialPageRoute(
                        builder: (context) => OrdersPages(
                          tester: widget.loginController.text,
                          title: 'Orders',
                        ),
                      ),
                    );
                  }
                },
              ),
            ],
          ),
        ),
      ),
    );
  }
}

List<String> parseUrl(String url) {
  return url.split('=');
}
