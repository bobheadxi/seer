import { Component, Prop, Mixins } from 'vue-property-decorator';
import { Pie } from 'vue-chartjs';

import { PlayerOverviews } from '@/math/aggregates';

@Component
export default class SimpleAggregatePie extends Mixins(Pie) {
  @Prop() aggs!: PlayerOverviews;
  @Prop() aggKeys!: string[];

  getDataByKey(): number[] {
    return Object.keys(this.aggs).map<number>((player) => {
      const playerAggs = this.aggs[player].aggs;
      if (!playerAggs) return 0;
      let v: any = playerAggs;
      this.aggKeys.forEach((k) => {
        // @ts-ignore
        v = v[k];
      });
      if (!v) return 0;
      return parseFloat(v);
    });
  }

  mounted() {
    this.renderChart({
      labels: Object.keys(this.aggs),
      datasets: [
        {
          backgroundColor: [
            'rgba(65, 184, 131, .8)',
            'rgba(228, 102, 81, .8)',
            'rgba(0, 216, 255, .8)',
            'rgba(155, 89, 182, .8)',
          ],
          borderWidth: 0,
          data: this.getDataByKey(),
        },
      ],
    }, {
      title: {
        text: 'Average Damage Dealt',
        display: true,
      },
      legend: {
        position: 'bottom',
      },
      responsive: true,
      maintainAspectRatio: false,
    });
  }
}
