#ifndef TETRISBOARD_H
#define TETRISBOARD_H

#include <QFrame>
#include <QBasicTimer>

// Классическая сетка Тетриса: 10 столбцов x 20 строк
const int BoardWidth = 10;
const int BoardHeight = 20;

class TetrisBoard : public QFrame
{
    Q_OBJECT

public:
    TetrisBoard(QWidget *parent = nullptr);

    void start();
    void pause();
    bool paused() const;

    QSize sizeHint() const override;

signals:
    void scoreChanged(int score);
    void linesRemovedChanged(int lines);
    void gameFinished();

protected:
    void paintEvent(QPaintEvent *event) override;
    void keyPressEvent(QKeyEvent *event) override;
    void timerEvent(QTimerEvent *event) override;

private:
    // Описание одной фигуры
    struct Shape {
        int coords[4][2]; // 4 клетки (x, y) относительно верхнего левого угла bounding box
        int color;        // индекс цвета
    };

    enum { NoShape = 0, NumShapes = 7, NumRotations = 4 };

    static const Shape shapes[NumShapes + 1][NumRotations];
    static const int colorTable[NumShapes + 1];

    void clearBoard();
    void newPiece();
    void drawSquare(QPainter &painter, int x, int y, int shapeIndex);
    bool tryMove(const Shape &shape, int newX, int newY);
    void rotatePiece();
    void pieceDropped();
    void hardDrop();
    void removeFullRows();
    void gameOver();

    // Состояние игры
    QBasicTimer timer;
    bool isStarted;
    bool isPaused;
    bool isGameOver;

    int nextPiece;
    int currentPiece;
    int curX;
    int curY;
    int curRot;
    int score;
    int linesRemoved;
    int numPiecesDropped;

    Shape curShape;
    int board[BoardWidth][BoardHeight];
};

#endif // TETRISBOARD_H