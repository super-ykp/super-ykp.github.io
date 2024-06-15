package enemy

import (
	"Neco/common"
	"Neco/danmaku"
	"math/rand"
)

// --=============================================================================
// ボス攻撃方法
// --=============================================================================
func stage1(t *TekiStruct, s bool) {
	if s {
		m_ZakoTamaSpeed = 1.5
		t.HpMax = 70
		return
	}
	if (t.flameCounter % 2) == 0 {
		danmaku.ShotB(m_dan, t.X, t.Y+16, 0.5, float64(int(t.flameCounter*10)%360), false, 2, 0.1)
	}
}

// -------------------------------------------------------------
func stage2(t *TekiStruct, s bool) {
	if s {
		m_ZakoTamaSpeed = 2.5
		t.HpMax = 70
		return
	}

	b_count1++
	if (b_count1 % 20) == 0 {
		danmaku.ShotB(m_dan, rand.Float64()*common.SCREEN_WIDTH, 0, 0.4, 0, true, 10, 0.1)
	}
}

// -------------------------------------------------------------
func stage3(t *TekiStruct, s bool) {
	if s {
		m_ZakoTamaSpeed = 0
		t.HpMax = 100
		return
	}

	if (t.flameCounter % 2) == 0 {
		danmaku.ShotB(m_dan, t.X, t.Y+16, 0.5, rand.Float64()*180-90, true, 10, 0.2)
		danmaku.ShotB(m_dan, t.X, t.Y+16, 0.5, rand.Float64()*180-90, true, 10, 0.4)
	}
}

// -------------------------------------------------------------
func stage4(t *TekiStruct, s bool) {
	if s {
		m_ZakoTamaSpeed = 0
		t.HpMax = 100
		return
	}
	r := float64((m_GlovalCounter * 5) % 160)
	sp := r/10 + 0.5
	if r <= 10 {
		danmaku.ShotB(m_dan, t.X, t.Y+16, sp, -rand.Float64()*20, true, 1, 0)
		danmaku.ShotB(m_dan, t.X, t.Y+16, sp, rand.Float64()*20, true, 1, 0)
	} else {
		danmaku.ShotB(m_dan, t.X, t.Y+16, sp*0.7, -r, true, 1, 0)
		danmaku.ShotB(m_dan, t.X, t.Y+16, sp*0.8, -r, true, 1, 0)
		danmaku.ShotB(m_dan, t.X, t.Y+16, sp*0.9, -r, true, 1, 0)
		danmaku.ShotB(m_dan, t.X, t.Y+16, sp, -r, true, 1, 0)
		danmaku.ShotB(m_dan, t.X, t.Y+16, sp*0.7, r, true, 1, 0)
		danmaku.ShotB(m_dan, t.X, t.Y+16, sp*0.8, r, true, 1, 0)
		danmaku.ShotB(m_dan, t.X, t.Y+16, sp*0.9, r, true, 1, 0)
		danmaku.ShotB(m_dan, t.X, t.Y+16, sp, r, true, 1, 0)
	}
}

// -------------------------------------------------------------
func stage5(t *TekiStruct, s bool) {
	if s {
		t.HpMax = 300
		m_ZakoTamaSpeed = 0.3
		return
	}

	if b_count2 == 0 {
		if b_count2 == 0 {
			//発射位置
			b_x = rand.Float64() * common.SCREEN_WIDTH
			b_y = 0
			//初期速度
			t.flameCounter = 0
			b_count2 = 1
		}
	} else if b_count2 == 1 {
		speed := float64(t.flameCounter)/10 + 2
		danmaku.ShotB(m_dan, b_x, b_y, speed, 0, true, 1, 0)
		if t.flameCounter > 30 {
			b_count2 = 0
		}
	}
}

// -------------------------------------------------------------
func stage6(t *TekiStruct, s bool) {
	if s {
		t.HpMax = 50
		m_ZakoTamaSpeed = 0
		b_count1 = 0
		return
	}
	b_count1++
	if b_count1 > common.SCREEN_WIDTH {
		b_count1 = 0
	}
	if (t.flameCounter % 5) == 0 {
		r := rand.Float64()*60 - 30
		danmaku.ShotB(m_dan, float64(b_count1), 0, 3, r, false, 1, 0)
	}
	d := rand.Float64()
	r := rand.Float64()*30 + 10
	if int(rand.Float64()*2)%2 == 0 {
		r = -r
	}
	danmaku.ShotB(m_dan, rand.Float64()*common.SCREEN_WIDTH/2-20, 0, 3+d, r, false, 1, 0)
	danmaku.ShotB(m_dan, rand.Float64()*common.SCREEN_WIDTH/2+20+common.SCREEN_WIDTH/2, 0, 3+d, r, false, 1, 0)
}

// -------------------------------------------------------------
func stage7(t *TekiStruct, s bool) {
	if s {
		t.HpMax = 100
		m_ZakoTamaSpeed = 0
		return
	}
	if (t.flameCounter % 30) == 0 {
		for i := 0; i < 360; i += 10 {
			r := t.flameCounter / 30
			danmaku.ShotB(m_dan, 0, common.SCREEN_HEIGHT/2, 0.5, float64(int(i+r)%360), true, 5, 0.1)
			danmaku.ShotB(m_dan, common.SCREEN_WIDTH, common.SCREEN_HEIGHT/2, 0.5, float64(int(i-r)%360), true, 5, 0.1)
		}
	}
}

// -------------------------------------------------------------
func stage8(t *TekiStruct, s bool) {
	if s {
		t.HpMax = 250
		m_ZakoTamaSpeed = 1
		return
	}
	m_GlovalCounter = 200
}

// -------------------------------------------------------------
func stage9(t *TekiStruct, s bool) {
	if s {
		t.HpMax = 50
		m_ZakoTamaSpeed = 0
		b_count2 = 0
		return
	}

	b_count1++
	b_count2 += (b_count1 / 50) * 180
	r := (b_count2) / 180
	if (t.flameCounter % 5) == 0 {
		for i := 0; i < 360; i += 60 {
			danmaku.ShotB(m_dan, t.X, t.Y, 3, float64(int(r+i)%360), true, 1, 0)
		}
	}
}

// -------------------------------------------------------------
func stage10(t *TekiStruct, s bool) {
	if s {
		t.HpMax = 250
		m_ZakoTamaSpeed = 0
		b_count1 = 100
		b_count2 = 0
		b_count3 = 0
		return
	}
	b_count2++
	if b_count2%(60*2) == 0 {
		b_count3 = ((b_count2 / 60 * 2) + 3 - 5) % 10
	}

	d := float64(b_count3)
	if (t.flameCounter%10) == 0 && b_count2%(60*2) > 12 {
		for i := 3; i < 10; i++ {
			danmaku.ShotB(m_dan, common.SCREEN_WIDTH/2-10+d, 0, 0.5, 0, true, 8, 0.2)
			danmaku.ShotB(m_dan, common.SCREEN_WIDTH/2-10+d+10, 0, 0.5, 0, true, 8, 0.2)
			danmaku.ShotB(m_dan, common.SCREEN_WIDTH/2-10+d+20, 0, 0.5, 0, true, 8, 0.2)
		}
	}

	b_count1++
	if b_count1 > 5*1000 {
		b_count1 = 5000
	}
	for i := 0; i < 10; i++ {
		sp := rand.ExpFloat64()*float64(b_count1)/1000 + 3
		danmaku.ShotB(m_dan, rand.Float64()*common.SCREEN_WIDTH/2-10, 0, sp, 0, true, 3, 0.2)
		danmaku.ShotB(m_dan, rand.Float64()*common.SCREEN_WIDTH/2+10+common.SCREEN_WIDTH/2, 0, sp, 0, true, 3, 0.2)
	}
}
