import gg

fn main() {
	gg.new_context(
		width: 800
		height: 600
		window_title: 'test-v'
		frame_fn: frame
	).run()
}

fn frame(mut ctx gg.Context) {
	ctx.begin()
	ctx.draw_text(100, 100, 'Hello from test-v (V Engine)', gg.TextCfg{
		size: 40
	})
	ctx.end()
}
