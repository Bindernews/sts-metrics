

function transformStatsOverview(data) {
    const rowNames = ['Name','Runs','Wins','Avg Win Rate', 'Deck Size (p25,p50,p75)','Floor Reached (p25,p50,p75)'];
    const order = ['name', 'runs', 'wins', 'avg_win_rate', 'p_deck_size', 'p_floor_reached'];
    const out = rowNames.map((n, i) => ({
        id: i,
        name: n,
    }));
    // Put each character into their own column
    data[0].forEach((row, row_i) => {
        const col = 'c' + row_i;
        for (let i = 0; i < order.length; i += 1) {
            out[i][col] = row[order[i]];
        }
    });
    return out;
}