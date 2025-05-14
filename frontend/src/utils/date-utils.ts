export enum Format {
    DayMonthYear = 'DD.MM.YYYY',
    DayMonthYearTime = 'DD.MM.YYYY HH:mm',
    MonthDayYear = 'MM.DD.YYYY',
    YearMonthDay = 'YYYY-MM-DD',
}

export function formatDate(date: Date, format: Format): string {
    const day = date.getDate().toString().padStart(2, '0');
    const month = (date.getMonth() + 1).toString().padStart(2, '0');
    const year = date.getFullYear();
    const hours = date.getHours().toString().padStart(2, '0');
    const minutes = date.getMinutes().toString().padStart(2, '0');

    switch (format) {
        case Format.DayMonthYear:
            return `${day}.${month}.${year}`;
        case Format.DayMonthYearTime:
            return `${day}.${month}.${year} ${hours}:${minutes}`;
        case Format.MonthDayYear:
            return `${month}.${day}.${year}`;
        case Format.YearMonthDay:
            return `${year}-${month}-${day}`;
        default:
            return date.toLocaleDateString();
    }
}
