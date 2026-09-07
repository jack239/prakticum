#ifndef TETRISWINDOW_H
#define TETRISWINDOW_H

#include <QWidget>

class QLCDNumber;
class QLabel;
class QPushButton;
class TetrisBoard;

class TetrisWindow : public QWidget
{
    Q_OBJECT

public:
    TetrisWindow(QWidget *parent = nullptr);

private slots:
    void startClicked();
    void pauseClicked();

private:
    void createLayout();

    TetrisBoard *board;
    QLCDNumber *scoreLcd;
    QLCDNumber *linesLcd;
    QLabel *gameOverLabel;
    QPushButton *startButton;
    QPushButton *pauseButton;
};

#endif // TETRISWINDOW_H