#include <QApplication>
#include "tetriswindow.h"

int main(int argc, char *argv[])
{
    QApplication app(argc, argv);

    TetrisWindow window;
    window.show();

    return app.exec();
}