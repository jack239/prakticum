const userName = null;          // имя участника, пока может быть не указано
const age = 19;                 // возраст участника
let categoryCode = 'half';      // код дистанции: 'fun', 'ten', 'half', 'full'
let hasMedicalCertificate = true; // есть ли медсправка
let promoCode;                  // промокод может быть не задан

let displayName = userName || 'Guest';
let distanceDescription;
switch (categoryCode) {
  case 'fun':
    distanceDescription = 'Фановый забег 5 км';
    break;
  case 'ten':
    distanceDescription = 'Спринт 10 км';
    break;
  case 'half':
    distanceDescription = 'Полумарафон 21 км';
    break;
  case 'full':
    distanceDescription = 'Марафон 42 км';
    break;
  default:
      distanceDescription = 'Неизвестная дистанция';
}


let isAllowed;

if (categoryCode === 'fun') {
    // фановый забег: возраст не ограничен, медсправка не обязательна
    isAllowed = true;
} else if (categoryCode === 'ten') {
    // 10 км: возраст от 16, медсправка обязательна
    if (age >= 16 && hasMedicalCertificate) {
    isAllowed = true;
    } else {
    isAllowed = false;
    }
} else if (categoryCode === 'half' || categoryCode === 'full') {
    // полумарафон и марафон: возраст от 18, медсправка обязательна
    if (age >= 18 && hasMedicalCertificate) {
    isAllowed = true;
    } else {
    isAllowed = false;
    }
} else {
    isAllowed = false;
}

const registrationStatus = isAllowed
    ? 'Допущен к забегу'
    : 'Не допущен к забегу';

const promoText = promoCode ?? 'Промокод не использован';

const report = `Участник: ${displayName}
    Возраст: ${age}
    Дистанция: ${distanceDescription}
    Статус регистрации: ${registrationStatus}
    Информация о промокоде: ${promoText}
`;

console.log(report);