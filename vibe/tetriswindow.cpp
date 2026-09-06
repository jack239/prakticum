#include "tetriswindow.h"
#include "tetrisboard.h"

#include <QLCDNumber>
#include <QLabel>
#include <QPushButton>
#include <QVBoxLayout>
#include <QHBoxLayout>
#include <QGroupBox>

TetrisWindow::TetrisWindow(QWidget *parent)
    : QWidget(parent)
{
    board = new TetrisBoard;
    board->setMinimumSize(board->sizeHint());

    scoreLcd = new QLCDNumber(5);
    scoreLcd->setSegmentStyle(QLCDNumber::Filled);
    scoreLcd->display(0);

    linesLcd = new QLCDNumber(5);
    linesLcd->setSegmentStyle(QLCDNumber::Filled);
    linesLcd->display(0);

    gameOverLabel = new QLabel(tr(""));
    gameOverLabel->setAlignment(Qt::AlignHCenter);
    gameOverLabel->setStyleSheet("color: #ff5050; font-weight: bold;");

    startButton = new QPushButton(tr("&Начать"));
    startButton->setFocusPolicy(Qt::NoFocus);
    pauseButton = new QPushButton(tr("&Пауза"));
    pauseButton->setFocusPolicy(Qt::NoFocus);
    pauseButton->setEnabled(false);

    createLayout();

    connect(startButton, &QPushButton::clicked, this, &TetrisWindow::startClicked);
    connect(pauseButton, &QPushButton::clicked, this, &TetrisWindow::pauseClicked);

    connect(board, &TetrisBoard::scoreChanged, scoreLcd, QOverload<int>::of(&QLCDNumber::display));
    connect(board, &TetrisBoard::linesRemovedChanged, linesLcd, QOverload<int>::of(&QLCDNumber::display));
    connect(board, &TetrisBoard::gameFinished, this, [this]() {
        gameOverLabel->setText(tr("Игра окончена! Нажмите «Начать» заново."));
        pauseButton->setEnabled(false);
        startButton->setEnabled(true);
        pauseButton->setText(tr("&Пауза"));
    });

    setWindowTitle(tr("Тетрис"));
    setFixedSize(sizeHint());
}

void TetrisWindow::createLayout()
{
    auto *rightPanel = new QGroupBox(tr("Информация"));
    auto *infoLayout = new QVBoxLayout;

    infoLayout->addWidget(new QLabel(tr("Очки")));
    infoLayout->addWidget(scoreLcd);
    infoLayout->addSpacing(15);
    infoLayout->addWidget(new QLabel(tr("Линии")));
    infoLayout->addWidget(linesLcd);

    rightPanel->setLayout(infoLayout);

    auto *buttonLayout = new QVBoxLayout;
    buttonLayout->addWidget(startButton);
    buttonLayout->addWidget(pauseButton);
    buttonLayout->addStretch(1);

    auto *rightLayout = new QVBoxLayout;
    rightLayout->addWidget(rightPanel);
    rightLayout->addLayout(buttonLayout);

    auto *boardLayout = new QHBoxLayout;
    boardLayout->addWidget(board);
    boardLayout->addLayout(rightLayout);

    auto *mainLayout = new QVBoxLayout;
    mainLayout->addLayout(boardLayout);
    mainLayout->addWidget(gameOverLabel);

    setLayout(mainLayout);
}

void TetrisWindow::startClicked()
{
    board->start();
    startButton->setEnabled(false);
    pauseButton->setEnabled(true);
    gameOverLabel->clear();
}

void TetrisWindow::pauseClicked()
{
    board->pause();
    pauseButton->setText(board->paused() ? tr("&Продолжить") : tr("&Пауза"));
}