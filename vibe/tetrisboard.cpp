#include "tetrisboard.h"

#include <QKeyEvent>
#include <QPainter>
#include <QRandomGenerator>

// Каждая фигура имеет 4 вращения, каждое вращение — 4 клетки (x, y).
const TetrisBoard::Shape TetrisBoard::shapes[NumShapes + 1][NumRotations] = {
    // NoShape (пусто)
    {
        { { {0,0},{0,0},{0,0},{0,0} },
          { {0,0},{0,0},{0,0},{0,0} },
          { {0,0},{0,0},{0,0},{0,0} },
          { {0,0},{0,0},{0,0},{0,0} } }
    },
    // I
    {
        { { {0,1},{1,1},{2,1},{3,1} },
          { {2,0},{2,1},{2,2},{2,3} },
          { {0,2},{1,2},{2,2},{3,2} },
          { {1,0},{1,1},{1,2},{1,3} } }
    },
    // O
    {
        { { {1,0},{2,0},{1,1},{2,1} },
          { {1,0},{2,0},{1,1},{2,1} },
          { {1,0},{2,0},{1,1},{2,1} },
          { {1,0},{2,0},{1,1},{2,1} } }
    },
    // T
    {
        { { {1,0},{0,1},{1,1},{2,1} },
          { {1,0},{1,1},{2,1},{1,2} },
          { {0,1},{1,1},{2,1},{1,2} },
          { {1,0},{0,1},{1,1},{1,2} } }
    },
    // S
    {
        { { {1,0},{2,0},{0,1},{1,1} },
          { {1,0},{1,1},{2,1},{2,2} },
          { {1,1},{2,1},{0,2},{1,2} },
          { {0,0},{0,1},{1,1},{1,2} } }
    },
    // Z
    {
        { { {0,0},{1,0},{1,1},{2,1} },
          { {2,0},{1,1},{2,1},{1,2} },
          { {0,1},{1,1},{1,2},{2,2} },
          { {1,0},{0,1},{1,1},{0,2} } }
    },
    // J
    {
        { { {0,0},{0,1},{1,1},{2,1} },
          { {1,0},{2,0},{1,1},{1,2} },
          { {0,1},{1,1},{2,1},{2,2} },
          { {1,0},{1,1},{0,2},{1,2} } }
    },
    // L
    {
        { { {2,0},{0,1},{1,1},{2,1} },
          { {1,0},{1,1},{1,2},{2,2} },
          { {0,1},{1,1},{2,1},{0,2} },
          { {0,0},{1,0},{1,1},{1,2} } }
    }
};

// Индексы NoShape..7 -> цвета (NoShape = 0 значит пусто/чёрный)
const int TetrisBoard::colorTable[NumShapes + 1] = {
    0, 1, 2, 3, 4, 5, 6, 7
};

TetrisBoard::TetrisBoard(QWidget *parent)
    : QFrame(parent)
{
    setFrameStyle(QFrame::Panel | QFrame::Sunken);
    setFocusPolicy(Qt::StrongFocus);
    isStarted = false;
    isPaused = false;
    isGameOver = false;
    nextPiece = NoShape;
    clearBoard();
}

void TetrisBoard::clearBoard()
{
    for (int i = 0; i < BoardHeight; ++i)
        for (int j = 0; j < BoardWidth; ++j)
            board[j][i] = NoShape;
}

void TetrisBoard::start()
{
    isStarted = true;
    isPaused = false;
    isGameOver = false;
    score = 0;
    linesRemoved = 0;
    numPiecesDropped = 0;

    clearBoard();
    emit scoreChanged(score);
    emit linesRemovedChanged(linesRemoved);

    nextPiece = NoShape;
    newPiece();
    timer.start(500, this);
}

bool TetrisBoard::paused() const
{
    return isPaused;
}

void TetrisBoard::pause()
{
    if (!isStarted || isGameOver)
        return;

    isPaused = !isPaused;
    if (isPaused)
        timer.stop();
    else
        timer.start(500, this);
    update();
}

QSize TetrisBoard::sizeHint() const
{
    return QSize(BoardWidth * 30 + frameWidth() * 2,
                 BoardHeight * 30 + frameWidth() * 2);
}

void TetrisBoard::paintEvent(QPaintEvent *)
{
    QPainter painter(this);
    QRect rect = contentsRect();

    painter.fillRect(rect, QColor(20, 20, 30));

    if (isGameOver) {
        painter.setPen(QColor(255, 80, 80));
        painter.setFont(QFont("Arial", 16, QFont::Bold));
        painter.drawText(rect, Qt::AlignCenter, tr("Игра окончена"));
        return;
    }

    if (isPaused) {
        painter.setPen(Qt::white);
        painter.setFont(QFont("Arial", 16, QFont::Bold));
        painter.drawText(rect, Qt::AlignCenter, tr("Пауза"));
        return;
    }

    const int cell = 30;

    // Нарисовать упавшие клетки
    for (int i = 0; i < BoardHeight; ++i) {
        for (int j = 0; j < BoardWidth; ++j) {
            int shape = board[j][i];
            if (shape != NoShape)
                drawSquare(painter, rect.left() + j * cell,
                           rect.top() + i * cell, shape);
        }
    }

    // Нарисовать текущую фигуру
    if (isStarted) {
        for (int i = 0; i < 4; ++i) {
            int x = curX + curShape.coords[i][0];
            int y = curY + curShape.coords[i][1];
            if (y >= 0)
                drawSquare(painter, rect.left() + x * cell,
                           rect.top() + y * cell, currentPiece);
        }
    }
}

void TetrisBoard::drawSquare(QPainter &painter, int x, int y, int shapeIndex)
{
    static const QColor colors[8] = {
        QColor(0, 0, 0),             // NoShape
        QColor(0, 200, 255),         // I  - голубой
        QColor(255, 220, 0),         // O  - жёлтый
        QColor(160, 0, 200),         // T  - фиолетовый
        QColor(0, 220, 0),           // S  - зелёный
        QColor(255, 60, 60),         // Z  - красный
        QColor(0, 120, 255),         // J  - синий
        QColor(255, 140, 0)          // L  - оранжевый
    };

    const int cell = 30;
    QColor color = colors[shapeIndex];
    QColor light = color.lighter(130);
    QColor dark = color.darker(150);

    painter.fillRect(x + 1, y + 1, cell - 2, cell - 2, color);

    painter.setPen(light);
    painter.drawLine(x, y + cell - 1, x, y);
    painter.drawLine(x, y, x + cell - 1, y);

    painter.setPen(dark);
    painter.drawLine(x + 1, y + cell - 1, x + cell - 1, y + cell - 1);
    painter.drawLine(x + cell - 1, y + cell - 1, x + cell - 1, y + 1);
}

void TetrisBoard::keyPressEvent(QKeyEvent *event)
{
    if (!isStarted || isGameOver) {
        QFrame::keyPressEvent(event);
        return;
    }

    if (isPaused) {
        if (event->key() == Qt::Key_P)
            pause();
        return;
    }

    switch (event->key()) {
    case Qt::Key_Left:
        tryMove(curShape, curX - 1, curY);
        break;
    case Qt::Key_Right:
        tryMove(curShape, curX + 1, curY);
        break;
    case Qt::Key_Down:
        tryMove(curShape, curX, curY + 1);
        break;
    case Qt::Key_Up:
        rotatePiece();
        break;
    case Qt::Key_Space:
        hardDrop();
        break;
    case Qt::Key_P:
        pause();
        break;
    default:
        QFrame::keyPressEvent(event);
    }
}

void TetrisBoard::rotatePiece()
{
    int newRot = (curRot + 1) % NumRotations;
    const Shape &newShape = shapes[currentPiece][newRot];
    if (tryMove(newShape, curX, curY))
        curRot = newRot;
}

void TetrisBoard::timerEvent(QTimerEvent *event)
{
    if (event->timerId() == timer.timerId()) {
        if (isStarted && !isPaused && !isGameOver) {
            if (!tryMove(curShape, curX, curY + 1))
                pieceDropped();
        }
    } else {
        QFrame::timerEvent(event);
    }
}

void TetrisBoard::newPiece()
{
    // Выбрать следующую фигуру случайно (1..7)
    nextPiece = QRandomGenerator::global()->bounded(NumShapes) + 1;
    currentPiece = nextPiece;
    curShape = shapes[currentPiece][0];
    curRot = 0;
    curX = BoardWidth / 2 - 2;
    curY = 0;

    if (!tryMove(curShape, curX, curY)) {
        // Фигура не может появиться — конец игры
        gameOver();
    }
}

bool TetrisBoard::tryMove(const Shape &shape, int newX, int newY)
{
    for (int i = 0; i < 4; ++i) {
        int x = newX + shape.coords[i][0];
        int y = newY + shape.coords[i][1];

        if (x < 0 || x >= BoardWidth || y >= BoardHeight)
            return false;
        if (y >= 0 && board[x][y] != NoShape)
            return false;
    }

    curShape = shape;
    curX = newX;
    curY = newY;
    update();
    return true;
}

void TetrisBoard::pieceDropped()
{
    // Зафиксировать текущую фигуру на поле
    for (int i = 0; i < 4; ++i) {
        int x = curX + curShape.coords[i][0];
        int y = curY + curShape.coords[i][1];
        board[x][y] = currentPiece;
    }

    ++numPiecesDropped;

    removeFullRows();

    // Немного ускоряемся с каждой фигурой
    timer.start(qMax(100, 500 - numPiecesDropped * 5), this);

    newPiece();
}

void TetrisBoard::hardDrop()
{
    while (tryMove(curShape, curX, curY + 1))
        ;
    pieceDropped();
}

void TetrisBoard::removeFullRows()
{
    int numFullRows = 0;
    int rowsToRemove[BoardHeight];
    int rowCount = 0;

    // Найти заполненные строки
    for (int y = BoardHeight - 1; y >= 0; --y) {
        bool full = true;
        for (int x = 0; x < BoardWidth; ++x) {
            if (board[x][y] == NoShape) {
                full = false;
                break;
            }
        }
        if (full) {
            rowsToRemove[rowCount++] = y;
            ++numFullRows;
        }
    }

    if (numFullRows == 0)
        return;

    // Сдвинуть строки вниз на количество удалённых строк выше
    int destRow = BoardHeight - 1;
    for (int srcRow = BoardHeight - 1; srcRow >= 0; --srcRow) {
        bool isRemoved = false;
        for (int i = 0; i < numFullRows; ++i) {
            if (rowsToRemove[i] == srcRow) {
                isRemoved = true;
                break;
            }
        }
        if (!isRemoved) {
            if (destRow != srcRow) {
                for (int x = 0; x < BoardWidth; ++x)
                    board[x][destRow] = board[x][srcRow];
            }
            --destRow;
        }
    }

    // Очистить верхние строки
    for (int y = 0; y <= destRow; ++y)
        for (int x = 0; x < BoardWidth; ++x)
            board[x][y] = NoShape;

    // Обновить счёт
    linesRemoved += numFullRows;
    // Очки: чем больше строк за раз, тем больше очков
    int points;
    switch (numFullRows) {
    case 1: points = 100; break;
    case 2: points = 300; break;
    case 3: points = 500; break;
    default: points = 800; break;
    }
    score += points;

    emit scoreChanged(score);
    emit linesRemovedChanged(linesRemoved);

    update();
}

void TetrisBoard::gameOver()
{
    isGameOver = true;
    timer.stop();
    emit gameFinished();
    update();
}